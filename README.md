## Get the environment ready FOR DEVELOPMENT

Certify yourself that the variable `GOPATH` is set and that `$GOPATH/bin` is in your `PATH`.

1. `go get github.com/mholt/caddy/caddy`
4. `go get github.com/upframe/fest`
5. Edit `GOPATH\src\github.com\mholt\caddy\caddyhttp\httpserver\plugin.go` and add `"upframe",` on line 434.
6. Edit `GOPATH\src\github.com\mholt\caddy\caddy\caddymain\run.go` and add `_ "github.com/upframe/fest"` to line 20.
5. `cd $GOPATH/src/github.com/upframe/fest/dist`
6. `go install github.com/mholt/caddy/caddy && caddy`

If you get any errors during the installation related to missing packages, run `go get` for those.

To build a final version, execute:

+ `go build -a -ldflags '-s' github.com/mholt/caddy/caddy -o caddy`
