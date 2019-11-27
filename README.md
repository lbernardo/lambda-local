# lambda-local
Execute lambda aws functions local

## Example

```bash
$ lambda-local start --volume $PWD
```

Start local lambda functions development


## Build
For build
```
make build
```
or
```
go build -o lambda-local
```

## Install
```
go install github.com/lbernardo/lambda-local
```

**Note: $GOPATH/bin must be set to $PATH**
