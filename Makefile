all: bins

clean:
	@echo cleaning...
	@GO111MODULE=on go clean -x
	@echo done!

bins:
	@echo building...
	@GO111MODULE=on go build ./...
	@echo done!

test: bins
	@echo testing...
	@GO111MODULE=on go test -covermode=count -coverprofile=coverage.out \
		github.com/freerware/negotiator/internal/representation \
		github.com/freerware/negotiator/internal/representation/xml \
		github.com/freerware/negotiator/internal/representation/json \
		github.com/freerware/negotiator/internal/representation/yaml \
		github.com/freerware/negotiator/internal/header \
		github.com/freerware/negotiator/proactive \
		github.com/freerware/negotiator/reactive \
		github.com/freerware/negotiator/transparent \
		github.com/freerware/negotiator/representation \
		github.com/freerware/negotiator \

	@echo done!

mocks:
	@echo mocking...
	@mockgen -source=representation/chooser.go -destination=internal/test/mock/chooser.go -package=mock -mock_names=Chooser=Chooser
	@echo done!
