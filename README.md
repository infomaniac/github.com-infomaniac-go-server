# Example Usage


This is the simplest way of using this package.
1. Define a fasthttp.RequestHandler
2. Create Server with this handler
3. Run the server

The `Run` method blocks until the server is stopped. This can either happen either via the `Stop()` method or via system signals, e.g. `SIGKILL` or `SIGINT`.  
Sending `SIGHUP` will reload the server by first stopping and the starting it again.


```go
func main() {
    s, err := server.NewGCP( handler() )
	if err != nil {
        log.Fatal(err)
	}
	s.Run()
}

func handler() fasthttp.RequestHandler {
    return func(ctx *fasthttp.RequestCtx) {
        ctx.Response.SetStatusCode(200)
	    ctx.Response.SetBodyString("Hello, World!")
    }
}
```