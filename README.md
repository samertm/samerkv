# samerkv

Distributed key-value database for the Bradfield School of Computer
Science distributed systems course (Fall 2020).

## Instructions

First, start the server in one shell:

```
cd ./cmd/server
go run .
```

Then, start a client in another shell:
```
cd ./cmd/client
go run .
```

You can issue the following commands:

```
get <key> # returns the value of the key, or the empty string if the key isn't set
set <key>=<value> # sets key to value
```


## TODO
 - Write tests
