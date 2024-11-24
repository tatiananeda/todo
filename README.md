## Study Project ToDo app api 

In memory storage

Implements CRUD on `\tasks` endpoint

Is ran on port `:3032`

See openapi file for details.

### to run benchmarking tests

`go test -bench=. ./controllers -run=^$`

### to run unit tests

`go test ./...  -v`