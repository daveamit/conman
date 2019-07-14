.PHONY: deps certs clean ca root etcd ci

CFSSL	= @env PATH=$(GOPATH)/bin:$(PATH) cfssl
JSON	= env PATH=$(GOPATH)/bin:$(PATH) cfssljson

deps:
	go get -u -tags nopkcs11 github.com/cloudflare/cfssl/cmd/cfssl
	go get -u github.com/cloudflare/cfssl/cmd/cfssljson
	go get -u github.com/mattn/goreman


certs: ca root etcd ci

clean:
	rm -f bundle.pem
	rm -rf certs

ca:
	mkdir -p certs
	$(CFSSL) gencert -initca config/ca/csr.json | $(JSON) -bare certs/ca

root:
	$(CFSSL) gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca-key.pem \
	  -config config/root/config.json \
	  config/root/csr.json | $(JSON) -bare certs/root
	  cat certs/ca.pem certs/root.pem > bundle.pem
	  mkdir -p certs/services
	  mkdir -p certs/users

etcd:
	$(CFSSL) gencert \
	  -ca certs/root.pem \
	  -ca-key certs/root-key.pem \
	  -config config/ca/config.json \
	  config/services/etcd/csr.json | $(JSON) -bare certs/services/etcd

ci:
	$(CFSSL) gencert \
	  -ca certs/root.pem \
	  -ca-key certs/root-key.pem \
	  -config config/ca/config.json \
	  config/users/ci/csr.json | $(JSON) -bare certs/users/ci
