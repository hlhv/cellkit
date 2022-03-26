# cellkit

Cellkit is a module that facilitates the creation of HLHV cells. In the future,
this module will contain many convenience functions to make creating a cell
easier.

## Client

The client package exposes basic functions for connecting to a queen cell and
handling incoming requests. Below is a basic program utilizing the client
functionality:

```
package main

import (
        "github.com/hlhv/protocol"
        "github.com/hlhv/testcell/client"
)

func main () {
        // create leash
        leash := client.NewLeash()
        leash.OnHTTP(onHTTP)

        // connect to server
        leash.Ensure (
                "localhost:2001",

                // mount on the root. @ refers to the default hostname, which
                // localhost is aliased to by default.
                []client.Mount {
                        { "@", "/" },
                },
                "",
        )
}

// http request handler
func onHTTP (
        band *client.Band,
        head *protocol.FrameHTTPReqHead,
) {
        // passing nil writes no headers
        band.WriteHTTPHead(200, nil)
        
        // write response body
        band.WriteHTTPBody([]byte("hello, world!"))
}
```

This program will respond to any incoming HTTP request with the text "hello
world!".
