package etcd_test

import (
	"conman"
	"conman/driver"
	. "conman/etcd"
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	t "github.com/coreos/etcd/pkg/transport"
	etcdCli "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/embed"
	"go.etcd.io/etcd/pkg/transport"
)

var dir = "default.etcd"
var endpoint = "https://localhost:2379"
var cli *etcdCli.Client
var etcd *Wrapper

func TestMain(m *testing.M) {
	e := setupEtcd(dir)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)
	var err error

	clientTLSConfig := &transport.TLSInfo{
		CAFile:        "../tools/tls/bundle.pem",
		CertFile:      "../tools/tls/certs/root.pem",
		TrustedCAFile: "../tools/tls/bundle.pem",
		KeyFile:       "../tools/tls/certs/root-key.pem",
		// ClientCertAuth: true,
	}

	clientConfig, err := clientTLSConfig.ClientConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cli, err = etcdCli.New(etcdCli.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: time.Second * 5,
		TLS:         clientConfig,
		// Username:    "root",
		// Password:    "week-password",
	})

	etcd, err = NewEtcd(&ConnDetails{
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

	acURL, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	cfg.ACUrls = []url.URL{*acURL}
	cfg.LCUrls = []url.URL{*acURL}

	// tlsInfo := transport.TLSInfo{}
	clientTLSConfig := t.TLSInfo{
		CAFile:         "../tools/tls/bundle.pem",
		CertFile:       "../tools/tls/certs/services/etcd.pem",
		KeyFile:        "../tools/tls/certs/services/etcd-key.pem",
		TrustedCAFile:  "../tools/tls/bundle.pem",
		ClientCertAuth: true,
	}

	cfg.ClientTLSInfo = clientTLSConfig

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

	_, err := NewEtcd(&ConnDetails{
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

func TestInitializeAndReset(t *testing.T) {
	ctx := context.Background()

	iv := &conman.InitializationVector{
		Secret: "secret",
	}
	err := etcd.Initialize(ctx, iv)
	if err != nil {
		t.Error("Failed to initialize", err)
	}

	users, err := cli.UserList(ctx)
	if err != nil {
		t.Error("Failed to retrive list of users", err)
	}

	foundRoot := false
	for _, user := range users.Users {
		if user == "root" {
			foundRoot = true
			break
		}
	}

	if !foundRoot {
		t.Error("Root user not found")
	}

	rsp, err := cli.Get(ctx, "$")
	if err != nil {
		t.Error("Failed to read $", err)
	}

	if rsp.Count == 0 {
		t.Error("Partially initialized, $ not found")
	}

	// Test reset
	err = etcd.Reset(ctx, false)
	if err != nil {
		t.Error("Failed to reset", err)
	}

	users, err = cli.UserList(ctx)
	if err != nil {
		t.Error("Failed to retrive list of users", err)
	}

	foundRoot = false
	for _, user := range users.Users {
		if user == "root" {
			foundRoot = true
			break
		}
	}

	if foundRoot {
		t.Error("Root user not cleared")
	}

	rsp, err = cli.Get(ctx, "$")
	if err != nil {
		t.Error("Failed to read $", err)
	}

	if rsp.Count != 0 {
		t.Error("$ not cleared")
	}

}

func TestInitializeErrorHandling(t *testing.T) {
	auth := newAuthMock()
	kv := newMockKV()
	etcd := NewFakeWrapper(kv, auth, nil)

	ctx := context.Background()
	iv := &conman.InitializationVector{}

	auth.roleAdd = errors.New("RoleAdd is disabled")
	err := etcd.Initialize(ctx, iv)
	if err != auth.roleAdd {
		t.Error("Must fail to initialize if roleAdd fails")
	}

	auth.roleAdd = nil
	auth.userAdd = errors.New("UserAdd is disabled")
	err = etcd.Initialize(ctx, iv)
	if err != auth.userAdd {
		t.Error("Must fail to initialize if UserAdd fails")
	}

	auth.userAdd = nil
	auth.userGrantRole = errors.New("UserGrantRole is disabled")
	err = etcd.Initialize(ctx, iv)
	if err != auth.userGrantRole {
		t.Error("Must fail to initialize if UserGrantRole fails")
	}

	auth.userGrantRole = nil
	kv.put = errors.New("Put is disabled")
	err = etcd.Initialize(ctx, iv)
	if err != kv.put {
		t.Error("Must fail to initialize if Put fails")
	}

	kv.put = nil
	auth.authEnable = errors.New("AuthEnable is disabled")
	err = etcd.Initialize(ctx, iv)
	if err != auth.authEnable {
		t.Error("Must fail to initialize if AuthEnable fails")
	}
}

func TestResetErrorHandling(t *testing.T) {
	auth := newAuthMock()
	kv := newMockKV()
	etcd := NewFakeWrapper(kv, auth, nil)

	ctx := context.Background()

	auth.authDisable = errors.New("Auth is disabled")
	err := etcd.Reset(ctx, false)
	if err != auth.authDisable {
		t.Error("Must fail to reset if AuthDisable fails")
	}

	auth.authDisable = nil
	kv.delete = errors.New("Put is disabled")
	err = etcd.Reset(ctx, false)
	if err != kv.delete {
		t.Error("Must fail to reset if delete fails")
	}

	kv.delete = nil
	auth.userRevokeRole = errors.New("UserRevokeRole is disabled")
	err = etcd.Reset(ctx, false)
	if err != auth.userRevokeRole {
		t.Error("Must fail to reset if userRevokeRole fails")
	}

	auth.userRevokeRole = nil
	auth.roleDelete = errors.New("RoleDelete is disabled")
	err = etcd.Reset(ctx, false)
	if err != auth.roleDelete {
		t.Error("Must fail to reset if roleDelete fails")
	}

	auth.roleDelete = nil
	auth.userDelete = errors.New("userDelete is disabled")
	err = etcd.Reset(ctx, false)
	if err != auth.userDelete {
		t.Error("Must fail to reset if userDelete fails")
	}
}
