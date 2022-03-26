package client

import (
        "github.com/hlhv/protocol"
)

/* OnHTTP specifies the http request handler function for the leash.
 */
func (leash *Leash) OnHTTP (
        callback func (band *Band, head *protocol.FrameHTTPReqHead),
) {
        leash.handles.onHTTP = callback
}
