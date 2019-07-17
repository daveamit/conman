package etcd

import (
	"crypto/tls"
	"strings"
	"time"

	etcdCli "go.etcd.io/etcd/clientv3"

	"go.etcd.io/etcd/pkg/transport"
)

// Wrapper is wrapper written over standard etcd-v3 client
type Wrapper struct {
	cn      *ConnDetails
	kv      etcdCli.KV
	auth    etcdCli.Auth
	watcher etcdCli.Watcher
}

// NewFakeWrapper make Wrapper for testing purpose
func NewFakeWrapper(kv etcdCli.KV, auth etcdCli.Auth, watcher etcdCli.Watcher) *Wrapper {
	return &Wrapper{
		kv:      kv,
		auth:    auth,
		watcher: watcher,
	}
}

// ConnDetails represents etcd connection details
type ConnDetails struct {
	Endpoint  string
	TLSConfig *transport.TLSInfo
}

// NewEtcd connects and initializes new etcd client wrapper
func NewEtcd(cn *ConnDetails) (*Wrapper, error) {
	var tlsConfig *tls.Config
	var err error

	if cn.TLSConfig != nil {
		tlsConfig, err = cn.TLSConfig.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	endpoints := strings.Split(cn.Endpoint, ",")
	cfg := etcdCli.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 5,
		TLS:         tlsConfig,
	}
	cli, err := etcdCli.New(cfg)

	ew := &Wrapper{}
	ew.kv = etcdCli.NewKV(cli)
	ew.watcher = etcdCli.NewWatcher(cli)
	ew.auth = etcdCli.Auth(cli)
	ew.cn = cn

	return ew, nil
}
