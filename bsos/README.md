
Building bsos 
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"
```