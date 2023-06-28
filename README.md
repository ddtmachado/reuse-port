# Testing

Start at least two server instances:

```bash
$ go run reuse.go
```

Hitting the root will return the server id:

```bash
$ curl localhost:9696
Server 319613c2-19b5-424b-8a5c-a1a0b5a8533b
$ curl localhost:9696
Server 93954fbe-6c67-49a9-a5cd-87a5fb4dd127
```

To simulate a long running task use the `/wait` path:

```bash
$ curl localhost:9696/wait?time=5s
Server 319613c2-19b5-424b-8a5c-a1a0b5a8533b waited for: 1s
```

You can use [vegeta](https://github.com/tsenart/vegeta) to validate no requests are lost during a server shutdown/restart:

```bash
$ echo "GET http://localhost:9696/" | vegeta attack -rate=10 -duration=15s | vegeta report
```

Sending a SIGTERM with `kill` to the process handling the requests, or just hitting `Ctrl+C` on the terminal it's running should stop it gracefully and the remaining requests will be automatically redirected to the remaining process.

You should be able to stop and start processes during a test without losing requests. Processes handling `/wait` requests will not receive new requests but wait the existing ones to complete before exiting.