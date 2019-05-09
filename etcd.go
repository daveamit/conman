package conman

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	etcd "go.etcd.io/etcd/clientv3"

	"go.etcd.io/etcd/pkg/transport"
)

// SettingUpdateKind is constant representing kind of update
type SettingUpdateKind int

const (
	// SettingKindUnknown indicates that current setting update kind is not known
	SettingKindUnknown SettingUpdateKind = 0 + iota

	// SettingKindInitial indicates that setting update was part fetched as part of initial pull
	SettingKindInitial

	// SettingKindPut indicates that setting update represents a created / updated operation
	SettingKindPut

	// SettingKindDeleted indicates that current setting was part of a delete operation
	SettingKindDeleted
)

// SettingUpdate represents a setting change that happened on etcd end
type SettingUpdate struct {
	Key         string
	RelativeKey string
	Value       string
	Tag         interface{}
	Kind        SettingUpdateKind
}

// EtcdWrapper is wrapper written over standard etcd-v3 client
type EtcdWrapper struct {
	kv      etcd.KV
	watcher etcd.Watcher
}

// EtcdConnDetails represents etcd connection details
type EtcdConnDetails struct {
	Endpoint  string
	TLSConfig *transport.TLSInfo
}

// NewEtcd connects and initializes new etcd client wrapper
func NewEtcd(e EtcdConnDetails) (*EtcdWrapper, error) {
	var tlsConfig *tls.Config
	var err error

	if e.TLSConfig != nil {
		tlsConfig, err = e.TLSConfig.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	cli, err := etcd.New(etcd.Config{
		Endpoints:   strings.Split(e.Endpoint, ","),
		DialTimeout: time.Second * 5,
		TLS:         tlsConfig,
	})

	ew := &EtcdWrapper{}
	ew.kv = etcd.NewKV(cli)
	ew.watcher = etcd.NewWatcher(cli)

	return ew, nil
}

// SettingUpdateQueueSize is the watch channel buffer size
const SettingUpdateQueueSize = 50

// Subscribe takes a prefix and watches under that prefix for any change, if somethings changes
// it will push that settingUpdate in the returned channel along with the tag. Tag is for book keeping.
// Cancling the context will unsubscribe, it is necessary to do so otherwise it may cause memory leaks
func (e *EtcdWrapper) Subscribe(ctx context.Context, prefix string, tag interface{}) (<-chan *SettingUpdate, error) {
	wc := e.watcher.Watch(ctx, prefix, etcd.WithPrefix())

	suChan := make(chan *SettingUpdate, SettingUpdateQueueSize)

	// Do initial Fetch
	rsp, err := e.kv.Get(ctx, prefix, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range rsp.Kvs {

		op := SettingKindInitial
		su := &SettingUpdate{
			Key: string(kv.Key),
			// if updated key is /a/b/c and substring under watch was a then
			// RelativeKey would be /b/c
			RelativeKey: strings.Replace(string(kv.Key), prefix, "", -1),
			Value:       string(kv.Value),
			Kind:        op,
			Tag:         tag,
		}

		suChan <- su
	}

	// Watch async for any future change
	go func() {
		// Read the channel till it closes
		for w := range wc {
			// Iterate though all the events and push them to appropriate channels
			for _, event := range w.Events {
				var op SettingUpdateKind

				if event.Type == etcd.EventTypePut {
					op = SettingKindPut
				} else if event.Type == etcd.EventTypeDelete {
					op = SettingKindDeleted
				}

				kv := event.Kv

				su := &SettingUpdate{
					Key: string(kv.Key),
					// if updated key is /a/b/c and substring under watch was a then
					// RelativeKey would be /b/c
					RelativeKey: strings.Replace(string(kv.Key), prefix, "", -1),
					Value:       string(kv.Value),
					Kind:        op,
					Tag:         tag,
				}
		
				suChan <- su
					}
		}
	}()

	return suChan, nil
}
