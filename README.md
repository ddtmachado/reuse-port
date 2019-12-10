# Testing

Start at least two servers in separate terminals:

```bash
$ go run reuse.go
```

Hitting the root will return the server id:

```bash
$ curl localhost:9696
{"server":"97a64fb6-7706-4aeb-b24c-f0afbdc5546c"}
$ curl localhost:9696
{"server":"97ef967d-c1ab-4e5a-b505-68b02142a7c1"}
```

To simulate a long running task use the `/wait` path:

```bash
$ curl localhost:9696/wait?time=5s
{"server":"97a64fb6-7706-4aeb-b24c-f0afbdc5546c","waitedFor":"5s"}
```

Use [vegeta](https://github.com/tsenart/vegeta) to test that no request is lost during a server process shutdown / restart:

```bash
$ echo "GET http://localhost:9696/" | vegeta attack -rate=10 -duration=15s | vegeta report
```

Sending a SIGTERM with `kill` to the process handling the requests, or just hitting `Ctrl+C` on the terminal it's running should stop it gracefully and the remaining requests will be automatically redirected to the remaining process.

You should be able to stop and start processes during a test without losing a single request. Processes handling `/wait` requests will not receive new requests but wait the existing ones to complete before exiting.