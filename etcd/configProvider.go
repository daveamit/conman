package etcd

import (
	"context"
	"strings"

	"conman"

	etcdCli "go.etcd.io/etcd/clientv3"
)

// Subscribe takes a prefix and watches under that prefix for any change, if somethings changes
// it will push that settingUpdate in the returned channel along with the tag. Tag is for book keeping.
// Cancling the context will unsubscribe, it is necessary to do so otherwise it may cause memory leaks
func (e *Wrapper) Subscribe(ctx context.Context, prefix string, tag interface{}) (<-chan *conman.SettingUpdate, error) {
	wc := e.watcher.Watch(ctx, prefix, etcdCli.WithPrefix())

	suChan := make(chan *conman.SettingUpdate, conman.SettingUpdateQueueSize)

	// Do initial Fetch
	rsp, err := e.kv.Get(ctx, prefix, etcdCli.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range rsp.Kvs {

		op := conman.SettingKindInitial
		su := &conman.SettingUpdate{
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
				var op conman.SettingUpdateKind

				if event.Type == etcdCli.EventTypePut {
					op = conman.SettingKindPut
				} else if event.Type == etcdCli.EventTypeDelete {
					op = conman.SettingKindDeleted
				}

				kv := event.Kv

				su := &conman.SettingUpdate{
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
