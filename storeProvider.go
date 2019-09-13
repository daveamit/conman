package conman

import "context"

// StoreProvider is abstraction over storage engine
type StoreProvider interface {
	// Initialize locks the storage and only allows authorized users
	// this might take time if its populated with data
	Initialize(ctx context.Context, cp CryptoProvider) error
	Reset(ctx context.Context, force bool) error

	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
}
