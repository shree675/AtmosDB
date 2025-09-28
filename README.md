# AtmosDB

A simple in-memory key-value store with concurrency and strict types.

## Demonstration

<img src="https://github.com/user-attachments/assets/5b128228-cc70-4218-a36d-6e9fddef58b7" width="600px" />

## About

A simple in-memory concurrent cache server and CLI, packaged into one.

Unlike Redis, AtmosDB doesn't infer datatypes but demands explicit types.

All operations are atomic and thread-safe.

AtmosDB tries to achieve a low memory footprint for large number of key-value pairs, with a tradeoff on latency for concurrent connections.

Server currently communicates on HTTP/1.1, making it possible to connect to the server through a driver written in any language (planned).

## Setup

Clone the repository and execute the following to install `atmos-cli` locally:

```powershell
$ go build
$ go install
```

Start the AtmosDB server:

```powershell
$ atmosdb
```

On a different terminal window, run the CLI to connect to the server:

```powershell
$ atmosdb cli <server_host_port>    # http://localhost:8080
```

> [!WARNING]  
> CLI does not run on Git Bash due to a [known issue](https://github.com/chzyer/readline/issues/191).

## Commands

These are the commands currently supported:

1. `db.version`  
   Prints the server's version.
2. `GET key`  
   Prints the value stored in the key along with its datatype, or _\<nil\>_.
3. `EXISTS key`  
   Prints _true_ if key exists, else _false_.
4. `SETINT key val`  
   Upserts the key-value pair, val must be integer.
5. `SETFLOAT key val`  
   Upserts the key-value pair, val must be float.
6. `SETSTR key val` or `SETSTR key "val with spaces"`  
   Upserts the key-value pair, val can be anything but will be stored as string.
7. `SETINT.TTL key val ttl`  
   Upserts the key-value pair with TTL, val must be integer, ttl is in seconds.
8. `SETFLOAT.TTL key val ttl`  
   Upserts the key-value pair with TTL, val must be float, ttl is in seconds.
9. `SETSTR.TTL key val ttl` or `SETSTR.TTL key "val with spaces" ttl`  
   Upserts the key-value pair with ttl, val can be anything but will be stored as string, ttl is in seconds.
10. `DEL key`  
    Deletes the key-value pair if exists.
11. `INCR key`  
    Increments the value by 1, stored value must be int.
12. `DECR key`  
    Decrements the value by 1, stored value must be int.

## Datatypes

**int**, **float** and **string** datatypes are currently supported.

## Version

v0.1

## Contributing

Contributions are welcome!

Improvements and TODOs:

- [ ] Implement better connection pooling, or look for other transport layer approaches for keeping client connections open.
- [ ] Support sending events and client subscriptions.
- [ ] Support arrays and sets and their corresponding operations.
- [ ] Support other basic commands like `INCRBY`, `DECRBY`, `EXPIREAT`, `GETEXP` or any other for common usages.
- [ ] Write drivers for languages like Java and Go.
- [ ] Implement transactions.
