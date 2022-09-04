build:
	go build

test:
	go test -v -race -cover -covermode=atomic -coverprofile=coverage.txt -coverpkg=github.com/ahsayde/yapl/yapl,github.com/ahsayde/yapl/internal/operator,github.com/ahsayde/yapl/internal/parser,github.com/ahsayde/yapl/internal/renderer ./...
	go tool cover -html=coverage.txt -o coverage.html
