
try:
	go run main.go \
	-imap-user="$$(teamvault-username --teamvault-config ~/.teamvault.json --teamvault-key=xwXyOn)" \
	-imap-password="$$(teamvault-password --teamvault-config ~/.teamvault.json --teamvault-key=xwXyOn)" \
	-imap-server="mail.benjamin-borbe.de:993" \
	-dry-run=true \
	-v=2

delete:
	go run main.go \
	-imap-user="$$(teamvault-username --teamvault-config ~/.teamvault.json --teamvault-key=xwXyOn)" \
	-imap-password="$$(teamvault-password --teamvault-config ~/.teamvault.json --teamvault-key=xwXyOn)" \
	-imap-server="mail.benjamin-borbe.de:993" \
	-dry-run=false \
	-v=2


precommit: ensure format generate test check
	@echo "ready to commit"

ensure:
	go mod verify
	go mod vendor

format:
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec go run -mod=vendor github.com/incu6us/goimports-reviser -project-name $$(go list -m) -file-path "{}" \;

generate:
	rm -rf mocks avro
	go generate -mod=vendor ./...

test:
	go test -mod=vendor -race $(shell go list -mod=vendor ./... | grep -v /vendor/)

check: lint vet errcheck

vet:
	go vet -mod=vendor $(shell go list -mod=vendor ./... | grep -v /vendor/)

lint:
	go run -mod=vendor golang.org/x/lint/golint -min_confidence 1 $(shell go list -mod=vendor ./... | grep -v /vendor/)

errcheck:
	go run -mod=vendor github.com/kisielk/errcheck -ignore '(Close|Write|Fprint)' $(shell go list -mod=vendor ./... | grep -v /vendor/)
