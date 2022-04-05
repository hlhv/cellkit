package client

import (
        "github.com/hlhv/scribe"
        "github.com/hlhv/protocol"
)

type ReqKind int

const (
        ReqKindMount   ReqKind = iota
        ReqKindUnmount
)

type Req interface {
        Kind() ReqKind
}

type ReqMount struct {
        promise chan error
        Host    string
        Path    string
}

type ReqUnmount struct {
        promise chan error
        Host    string
        Path    string
}

func (req *ReqMount)   Kind () ReqKind { return ReqKindMount   }
func (req *ReqUnmount) Kind () ReqKind { return ReqKindUnmount }

func (leash *Leash) addQueue (req Req) {
        leash.queue <- req
}

/* Mount tells the leash to mount on a particular pattern. This function is
 * thread safe.
 */
func (leash *Leash) Mount (host string, path string) (err error) {
        scribe.PrintProgress(scribe.LogLevelNormal, "mounting on", host, path)
        promise := make(chan error)
        leash.addQueue (&ReqMount {
                promise: promise,
                Host:    host,
                Path:    path,
        })
        return <- promise
}

/* Unmount tells the leash to unmount off of a particular pattern. This function
 * is thread safe.
 */
func (leash *Leash) Unmount (host string, path string) (err error) {
        scribe.PrintProgress (
                scribe.LogLevelNormal, "unmounting from", host, path)
        promise := make(chan error)
        leash.addQueue (&ReqUnmount {
                promise: promise,
                Host:    host,
                Path:    path,
        })
        return <- promise
}

func (leash *Leash) respond () {
        for {
                req := <- leash.queue
                scribe.PrintRequest (
                        scribe.LogLevelDebug, "got internal request")
                if req == nil { break }
                leash.respondOnce(req)
        }
        scribe.PrintWarning(scribe.LogLevelDebug, "will no longer respond")
}

func (leash *Leash) respondOnce (req Req) {
        switch req.Kind() {
        case ReqKindMount:
                reqSure := req.(*ReqMount)
                _, err := leash.writeMarshalFrame (&protocol.FrameMount {
                        Host: reqSure.Host,
                        Path: reqSure.Path,
                })
                reqSure.promise <- err
                break
                
        case ReqKindUnmount:
                reqSure := req.(*ReqUnmount)
                _, err := leash.writeMarshalFrame (&protocol.FrameUnmount {
                        Host: reqSure.Host,
                        Path: reqSure.Path,
                })
                reqSure.promise <- err
                break
        }
}
