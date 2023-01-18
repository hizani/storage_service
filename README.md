# Microservice data storage application
There are two services: [http_server](http_server) and [storage_service](storage_service)

## storage_service
Storage gRPC service. 
Saves data in a storage of choice: *runtime*, *file*, *database*
## Usage
1. Set up `config.toml` file
2. Run `go run cmd/main.go STORAGE CONFIG` where **STORAGE** is one of the following arguments: `runtime`, `file`, `database` and **CONFIG** is a path to the config file.
### Example
```
go run cmd/storage_service/main.go runtime config.toml
```

---

## http_server
HTTP server that communicate with [storage](storage) via gRPC. 
### Usage
Run `go run cmd/http_server/main.go` with two sockets: http server socket and socket of the storage service
### Example
```
go run cmd/http_server/main.go localhost:8080 localhost:9786
```
Then try to send HTTP requests. Examples provided [here](requests.httpbook)
