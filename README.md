# random-cid

Generate a random CID

## Building
```bash
git clone https://github.com/application-research/random-cid.git
cd random-cid
go build .
```

### Usage
You can use the CLI directly
```bash
# random v1 CID
$ ./random-cid
bafkreigd6pc65n4mcfabl3ucv2waoxwizxqizrd32r6cmqjx5zsul4jv3i

# random v0 CID
$ ./random-cid -c 0
QmaA14Co9Q9AuNHcs6KH2ZmJ8sCTwW6ZN7TJfxNcXnrUAX
```

Or you can start an API to retrieve CIDs over HTTP
```bash
./random-cid api
```

The API endpoints are:
- `GET /`: returns a random CID v1
- `GET /v0`: returns a random CID v0
- `GET /v1`: returns a random CID v1

