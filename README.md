## Get the environment ready

Certify yourself that the variable `GOPATH` is set and that `$GOPATH/bin` is in your `PATH`.

1. `go get github.com/mholt/caddy/caddy`
2. `go get github.com/caddyserver/caddydev`
3. `go install github.com/caddyserver/caddydev`
4. `go get github.com/hacdias/upframe`
5. `cd $GOPATH/src/github.com/hacdias/upframe/dist`
6. `caddydev --source="github.com/hacdias/upframe" upframe`

If you get any errors during the installation related to missing packages, run `go get` for those.
