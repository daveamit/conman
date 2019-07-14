package conman

import "context"

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

// SettingUpdateQueueSize is the watch channel buffer size
const SettingUpdateQueueSize = 50

// ConfigProvider exposes a set of functionality that are expected out of
// a typical configuration provider
type ConfigProvider interface {
	Subscribe(ctx context.Context, prefix string, tag interface{}) (<-chan *SettingUpdate, error)
}
