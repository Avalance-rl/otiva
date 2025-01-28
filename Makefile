gofmt:
	gofumpt -l -w .
	goimports -w .
.PHONY:gofmt