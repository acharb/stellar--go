test: go-master
	CGO_ENABLED=1 go test -race -tags=refactortests

go-master:
	git clone --depth=1 https://github.com/stellar/go go-master
	cd go-master && \
		go mod edit -module=github.com/stellar/go-master
	go mod edit \
		-require=github.com/stellar/go-master@v0.0.0-00010101000000-000000000000 \
		-replace=github.com/stellar/go-master=./keypair/refactortests/go-master
