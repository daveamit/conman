package etcd_test

import (
	"conman"
	"conman/driver"
	. "conman/etcd"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	etcdCli "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/embed"
	"go.etcd.io/etcd/pkg/transport"
)

var dir = "default.etcd"
var endpoint = "http://127.0.0.1:2379"
var cli *etcdCli.Client
var etcd *Wrapper

func TestMain(m *testing.M) {
	e := setupEtcd(dir)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)
	var err error

	serverTLSConfig := &transport.TLSInfo{
		CAFile:   "../tools/tls/bundle.pem",
		CertFile: "../tools/tls/certs/services/etcd.pem",
		KeyFile:  "../tools/tls/certs/services/etcd-key.pem",
	}
	serverTlsInfo, err := serverTLSConfig.ClientConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cli, err = etcdCli.New(etcdCli.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: time.Second * 5,
		TLS:         serverTlsInfo,
	})

	clientTLSConfig := &transport.TLSInfo{
		CAFile:   "../tools/tls/bundle.pem",
		CertFile: "../tools/tls/certs/users/ci.pem",
		KeyFile:  "../tools/tls/certs/users/ci-key.pem",
	}

	etcd, err = NewEtcd(ConnDetails{
		Endpoint:  endpoint,
		TLSConfig: clientTLSConfig,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print out all the KVs

	fmt.Println(KVs)

	code := m.Run()
	e.Close()
	os.Exit(code)
}

func setupEtcd(dir string) *embed.Etcd {
	cfg := embed.NewConfig()
	cfg.Dir = dir
	// tlsInfo := transport.TLSInfo{}

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}
	// log.Fatal(<-e.Err())
	return e
}
func TestUnsubscribeOnContextCancel(t *testing.T) {
	ctx, done := context.WithCancel(context.Background())
	done()
	_, err := etcd.Subscribe(ctx, "section2", nil)
	if err == nil {
		t.Error("Subscribe must fail is context is cancel")
	}
}

func TestWrongTLSInfoMustFail(t *testing.T) {
	clientTLSConfig := &transport.TLSInfo{
		CAFile:   "../tools/tls/invalid.pem",
		CertFile: "../tools/tls/certs/users/invalid.pem",
		KeyFile:  "../tools/tls/certs/users/invalid.pem",
	}

	_, err := NewEtcd(ConnDetails{
		Endpoint:  endpoint,
		TLSConfig: clientTLSConfig,
	})
	if err == nil {
		t.Error("Must fail if invalid tlsInfo")
	}
}

func assertUpdate(t *testing.T, update *conman.SettingUpdate) {
	found := false

	for _, entry := range KVs {
		if entry.key == update.Key && entry.kind == update.Kind {
			if update.Kind == conman.SettingKindInitial || update.Kind == conman.SettingKindPut {
				if entry.value == update.Value {
					found = true
					break
				}
				continue
			}
			found = true
			break
		}
	}
	if !found {
		t.Logf("---------------------%s-%d----------------------------", update.Key, update.Kind)
		t.Log("update: ", update)
		t.Error("Unknown setting update")
		t.Log("------------------------------------------------------")
	}

}

type testKV struct {
	kind  conman.SettingUpdateKind
	key   string
	value string
}

var KVs = []testKV{
	testKV{conman.SettingKindInitial, "section1/section2/intact", "intact-value-will-never-change"},
	testKV{conman.SettingKindInitial, "section1/section2/section3", "first-value-must-update"},
	testKV{conman.SettingKindPut, "section1/section2/section3", "second-value-must-be-removed"},
	testKV{conman.SettingKindDeleted, "section1/section2/section3", "second-value-must-be-removed"},
}

func TestInitialUpdateAndDeleteWorks(t *testing.T) {
	ctx, done := context.WithCancel(context.Background())

	initialUpdates := 0
	// Seed the etcd, to see if we receive initial update
	for _, kv := range KVs {
		if kv.kind == conman.SettingKindInitial {
			cli.KV.Put(ctx, kv.key, kv.value)
			initialUpdates++
		}
	}

	// All initial updates will be fetched in one shot
	wg := sync.WaitGroup{}
	wg.Add(initialUpdates)

	updateStream, err := etcd.Subscribe(ctx, "section1", nil)
	go func() {
		t.Log("In go func")
		if err != nil {
			t.Error("Failed to subscribe", err)
			return
		}

		t.Log("Waiting for settings")
		for update := range updateStream {
			assertUpdate(t, update)
			wg.Done()
		}
	}()
	wg.Wait()
	wg.Add(len(KVs) - initialUpdates)

	for _, kv := range KVs {
		if kv.kind == conman.SettingKindPut {
			t.Logf("Putting - %s", kv.key)
			cli.KV.Put(ctx, kv.key, kv.value)
		}

		if kv.kind == conman.SettingKindDeleted {
			t.Logf("Removing - %s", kv.key)
			cli.KV.Delete(ctx, kv.key)
		}
	}
	wg.Wait()
	done()
}

func TestAssertDriverCompatibility(t *testing.T) {
	d := driver.New(etcd)

	wg := sync.WaitGroup{}
	wg.Add(1)

	updateHandler := func(update *conman.SettingUpdate) {
		assertUpdate(t, update)
		wg.Done()
	}

	err := d.Watch("section1", updateHandler, nil)
	if err != nil {
		t.Error("Failed to watch", err)
	}

	wg.Wait()
}
