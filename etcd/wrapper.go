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
	kv      etcdCli.KV
	auth    etcdCli.Auth
	watcher etcdCli.Watcher
}

// ConnDetails represents etcd connection details
type ConnDetails struct {
	Endpoint  string
	TLSConfig *transport.TLSInfo
}

// NewEtcd connects and initializes new etcd client wrapper
func NewEtcd(e ConnDetails) (*Wrapper, error) {
	var tlsConfig *tls.Config
	var err error

	if e.TLSConfig != nil {
		tlsConfig, err = e.TLSConfig.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	cli, err := etcdCli.New(etcdCli.Config{
		Endpoints:   strings.Split(e.Endpoint, ","),
		DialTimeout: time.Second * 5,
		TLS:         tlsConfig,
	})

	ew := &Wrapper{}
	ew.kv = etcdCli.NewKV(cli)
	ew.watcher = etcdCli.NewWatcher(cli)
	ew.auth = etcdCli.Auth(cli)

	return ew, nil
}
