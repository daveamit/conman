package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/daveamit/conman"
)

// ConfigProvider is
type ConfigProvider struct {
	provider conman.ConfigProvider

	// this his prefix is used to watch for dynamic creation / deletion
	// of buckets
	bucketListPrefix string
	buckets          []string
	mx               sync.Mutex

	// keys subscribed under specific buckets
	bucketedKeys sync.Map

	// keys currently being watched
	watchList sync.Map
}

// SettingChangeHandler represents a call back which is invoked when a setting is changed
type SettingChangeHandler func(update *conman.SettingUpdate)

// Watch should be used to watch under a static or fixed key
func (cp *ConfigProvider) Watch(setting string, handler SettingChangeHandler) error {
	_, found := cp.watchList.Load(setting)
	// If found any previously registered setting with same key, return error
	if found {
		return conman.ErrAlreadyWatchingGivenKey
	}

	ctx, cancelContext := context.WithCancel(context.Background())

	updates, err := cp.provider.Subscribe(ctx, setting, nil)
	if err != nil {
		// cancel the context if there was an error while subscribing
		cancelContext()
		return err
	}

	// Watch for setting updates and invoke hander
	go func() {
		for update := range updates {
			handler(update)
		}
	}()

	// add to our watchlist
	cp.watchList.Store(setting, cancelContext)

	return nil
}

// WatchBucketWise should be used to watch */key/* kind of settings
func (cp *ConfigProvider) WatchBucketWise(setting string, handler SettingChangeHandler) {
	cp.bucketedKeys.Store(setting, handler)

	cp.syncWatchList()
}

func (cp *ConfigProvider) syncWatchList() {
	for _, bucket := range cp.buckets {

		cp.bucketedKeys.Range(func(bucketedKey interface{}, handler interface{}) bool {
			// if bucket is bravo, bucketedKey is echo
			// key would be bravo/key
			key := fmt.Sprintf("%s/%s", bucket, bucketedKey)

			_, found := cp.watchList.Load(key)
			// if not found in watchlist, remove from watch list
			if !found {
				// Ignoring errors
				cp.Watch(key, handler.(SettingChangeHandler))
			}

			return true
		})
	}
}

// New initializes a configuration provider, this wrapper is designed
// for a programattic access prospective
func New(cp conman.ConfigProvider) *ConfigProvider {
	p := &ConfigProvider{
		provider:         cp,
		bucketListPrefix: "bucket",
	}

	bucketHandler := func(update *conman.SettingUpdate) {
		switch update.Kind {
		case conman.SettingKindInitial:
			fallthrough
		case conman.SettingKindPut:
			if !contains(p.buckets, update.RelativeKey) {
				p.mx.Lock()
				p.buckets = append(p.buckets, update.RelativeKey)
				p.syncWatchList()
				p.mx.Unlock()
			}
			// TODO: Bucket removed logic
		}
	}

	// Watch for dynamic buckets
	p.Watch(p.bucketListPrefix, bucketHandler)
	return p
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
