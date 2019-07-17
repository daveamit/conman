package conman

import "context"

// InitializationVector is used to initialize the storage driver
type InitializationVector struct {
	// Secret is required to encrypt data at reset, it it's empty, crypto is not enabled
	Secret string
}

// StoreProvider is abstraction over storage engine
type StoreProvider interface {
	// Initialize locks the storage and only allows authorized users
	// this might take time if its populated with data
	Initialize(ctx context.Context, iv InitializationVector) error
	Reset(ctx context.Context, force bool) error

	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
}
