CMD_GO=go

###
### test targets
###
.PHONY: test
test:
	$(CMD_GO) test -race -cover ./...

.PHONY: vet
vet:
	$(CMD_GO) vet ./...

.PHONY: vet-gen
vet-gen:
	$(CMD_GO) vet generator/main.go

.PHONY: generate
generate:
	$(CMD_GO) generate