# Intern project for wb
Storage service which saves data in a storage of choice: *runtime*, *file*, *database*
## Usage
1. Set up `config.toml` file
2. Run `go run cmd/main.go` with one of the following arguments: `runtime`, `file`, `database`
### Example
```
go run cmd/main.go runtime
```