# samerkv

Distributed key-value database for the Bradfield School of Computer
Science distributed systems course (Fall 2020).

## Instructions

First, start the synchronous follower replica in one shell:

```
go run ./cmd/server -type=syncfollower
```

Then, start the asynchronous follower replica in another shell:

```
go run ./cmd/server -type=asyncfollower
```

Then, start the leader replica in another shell:

```
go run ./cmd/server -type=leader
```

Finally, start a client in another shell:
```
go run ./cmd/client
```

You can issue the following commands:

```
get <key> # returns the value of the key, or the empty string if the key isn't set. uses the default table.
get <table>.<key> # returns the value in a given table

set <key>=<value> # sets key to value
set <table>.<key>=<value> # sets key to value in a given table
set *.<key>=<value> # sets key to a value in all tables

create_table <table> # creates a table
delete_table <table> # deletes a table
list_tables # lists tables
```

Every table comes with a default table, named "default". All get and
set requests that don't specify a table operate on the default table.
It's possible to delete the default table, but it's not possible to
set a new table as the default.

## Hacking

You need some grpc tooling installed locally to generate the protobufs.

```
brew install protoc
go get github.com/golang/protobuf/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

See the [grpc Go
quickstart](https://grpc.io/docs/languages/go/quickstart/) for more
information.

Generate new protobuf scaffolding with `make gen`.


## TODO
 - Write tests
