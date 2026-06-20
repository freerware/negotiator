all: bins

clean:
	@echo cleaning...
	@go clean -x
	@echo done!

bins:
	@echo building...
	@go build ./...
	@echo done!

test: bins
	@echo testing...
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo done!

lint:
	@echo linting...
	@golangci-lint run
	@echo done!

tidy:
	@echo tidying...
	@go mod tidy
	@echo done!

bench:
	@echo benchmarking...
	@go test -bench=. -benchmem ./...
	@echo done!

mocks:
	@echo mocking...
	@mockgen -source=representation/chooser.go -destination=internal/test/mock/chooser.go -package=mock -mock_names=Chooser=Chooser
	@echo done!
