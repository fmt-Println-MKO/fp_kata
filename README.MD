



## generate the coverage report
```shell
go test -cover -coverpkg=$(go list ./... | grep -v '/mocks' | grep -v '/cmd' | tr '\n' ',') ./... -coverprofile=coverage.out
```
