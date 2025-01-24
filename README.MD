
## Demo app in non-functional style




## generate dependency injection trough wire
```shell
cd internal/app && go generate ./
```

## generate mockerey mocks
```shell
mockery --all --output=mocks
```

## generate the coverage report
```shell
go test -cover -coverpkg=$(go list ./... | grep -v '/mocks' | grep -v '/cmd' | tr '\n' ',') ./... -coverprofile=coverage.out
```

