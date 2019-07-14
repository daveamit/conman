
deps:
	go get github.com/jpoles1/gopherbadger
	cd tools/tls && $(MAKE) deps

test:
	cd tools/tls && $(MAKE) certs
	go test --mod=vendor ./... && \
	$(MAKE) clean

cover:
	cd tools/tls && $(MAKE) certs
	go test --mod=vendor -coverprofile=c.out ./... && \
	go tool cover -html=c.out && \
	gopherbadger -md="README.md"
	cd tools/tls && $(MAKE) clean


clean:
	rm -f c.out
	cd tools/tls && $(MAKE) clean