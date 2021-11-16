target: ;

test:
	@go test ./pkg/parser

test_without_cache:
	@go clean -testcache && make test