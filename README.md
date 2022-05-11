# cellkit

---

!!! This repository is no longer in development. For the new cell module, visit
[this repository](https://github.com/hlhv/cell) !!!

---

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
        "github.com/hlhv/cellkit/frame"
        "github.com/hlhv/cellkit/client"
)

func main () {
        // configure framework
        frame.Be (&frame.Conf {
                Description: "Test cell",
                Run:         run,
        })
}

// main cell function
func run (leash *client.Leash) {
        leash.OnHTTP(onHTTP)
        
        // connect to server
        leash.Ensure (
                "localhost:2001",

                // mount on the root. @ refers to the default hostname, which
                // localhost is aliased to by default.
                client.Mount { "@", "/" },
                "something",
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
