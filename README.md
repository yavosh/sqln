# SQLITE Network 

A network accessible sqlite database using grpc protocol. Used to experiment with k8s operators.
The rqlite project already implements most of the goals of this project.

## Running 

For an im-memory database:

    go run cmd/sqln/main.go --data-source ":memory:"

## Examples: Querying

    grpcurl -plaintext -d '{"query": "SELECT * FROM tab"}' localhost:5051 sqln.Query/ExecuteQuery

## Notes

- See rqlite for a better sqlite server implementation :)
- See https://github.com/bsm/redeo/tree/master/resp for redis protocol interface
