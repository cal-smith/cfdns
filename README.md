# cfdns
## Utility to update Cloudflare DNS records with the current ip address

### Usage

```
$ cfdns -h
Usage for cfdns:
  -token string
        Cloudflare API token. Must have edit access to all zone:domain pairs specified
  -y    Set this flag to skip confirmation of changes
  [zone:domain]
        The domains to update are provided as a comma separated list of zone:domian pairs
        For example: foo.com:cloud,bar.net:www
```

A simple cron example:
```
# every hour update the ip address for the "cloud" subdomain
# logs output to the system logger with the 'cfdns' tag
0 * * * * /usr/local/bin/cfdns -token [token] -y example.com:cloud 2>&1 | logger -t cfdns
```

### Building

Most recently tested with go 1.18.1 - `go get && go build`

Directly depends on [cloudflare-go](https://github.com/cloudflare/cloudflare-go) and [term](https://pkg.go.dev/golang.org/x/term) 
