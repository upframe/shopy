## Get the environment ready

Certify yourself that the variable `GOPATH` is set and that `$GOPATH/bin` is in your `PATH`.

1. `go get github.com/mholt/caddy/caddy`
4. `go get github.com/hacdias/upframe`
5. Edit `GOPATH\src\github.com\mholt\caddy\caddyhttp\httpserver\plugin.go` and add `"upframe",` on line 434.
6. Edit `GOPATH\src\github.com\mholt\caddy\caddy\caddymain\run.go` and add `_ "github.com/hacdias/upframe"` to line 20.
5. `cd $GOPATH/src/github.com/hacdias/upframe/dist`
6. `go install github.com/mholt/caddy/caddy && caddy`

If you get any errors during the installation related to missing packages, run `go get` for those.
