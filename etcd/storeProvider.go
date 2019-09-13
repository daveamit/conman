package etcd

import (
	"conman"
	"context"
	"errors"

	"go.uber.org/zap"
)

// Initialize locks the storage and only allows authorized users
// this might take time if its populated with data
// 1. If $ is already there return with error
// 2. Create root user
// 3. Cache crypto Provider
func (w *Wrapper) Initialize(ctx context.Context, crypto conman.CryptoProvider) error {
	log := conman.Log().With(zap.String("@", "Initialize"))
	_, err := w.auth.RoleAdd(ctx, "root")
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to create root role")
		return err
	}

	// TODO: Replate this with random password
	_, err = w.auth.UserAdd(ctx, "root", "week-password")
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to create root user")
		return err
	}

	_, err = w.auth.UserGrantRole(ctx, "root", "root")
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to grant root role to root")
		return err
	}

	_, err = w.kv.Put(ctx, "$", "true")
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to put $")
		return err
	}

	_, err = w.auth.AuthEnable(ctx)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to lock")
		return err
	}

	// crypto is cached so that it can be used to read / write to encrypt / decrypt
	w.crypto = crypto

	return nil
}

// Reset unlocks and decrypts data at reset.
// 1. If $ not found, return with error
// 2. Decrypt every key
// 3. remove $
func (w *Wrapper) Reset(ctx context.Context, force bool) error {
	log := conman.Log().With(zap.String("@", "Reset"))

	log.Debug("Disabling auth")
	_, err := w.auth.AuthDisable(ctx)
	if err != nil {
		return err
	}

	// TODO: Decrypt keys

	_, err = w.kv.Delete(ctx, "$")
	if err != nil {
		return err
	}

	_, err = w.auth.UserRevokeRole(ctx, "root", "root")
	if err != nil {
		return err
	}

	_, err = w.auth.RoleDelete(ctx, "root")
	if err != nil {
		return err
	}

	_, err = w.auth.UserDelete(ctx, "root")
	if err != nil {
		return err
	}
	return nil
}

// Get reads a key, it decrypts it if encrypted
func (w *Wrapper) Get(ctx context.Context, key string) (string, error) {

	rsp, err := w.kv.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if rsp.Count == 0 {
		return "", errors.New("Key not found")
	}

	value := rsp.Kvs[0].Value
	if w.crypto != nil {
		value, err = w.crypto.Decrypt(value)
		if err != nil {
			return "", err
		}
	}

	return string(value), nil
}

// Set sets value for given key, encrypts if secret is provided
func (w *Wrapper) Set(ctx context.Context, key string, val string) error {
	var err error
	value := []byte(val)
	if w.crypto != nil {
		value, err = w.crypto.Encrypt(value)
		if err != nil {
			return err
		}
	}

	w.kv.Put(ctx, key, string(value))
	return nil
}
