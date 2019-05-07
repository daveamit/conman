package conman

import "context"

// ConfigProvider exposes a set of functionality that are expected out of
// a typical configuration provider
type ConfigProvider interface {
	Subscribe(ctx context.Context, prefix string, tag interface{}) (<-chan *SettingUpdate, error)
}
