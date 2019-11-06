go test ./... -coverprofile cover.out.tmp
cat cover.out.tmp | grep -v "_mock.go" > cover.out
go tool cover -func cover.out