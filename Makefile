
tidy:
	find . -name "*.go" -type f -not -path "./vendor/*" | xargs -n1 go fmt
	go mod tidy

test:
	cd cert && go test -v
	cd log && go test -v 

