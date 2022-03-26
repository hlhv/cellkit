package client

import (
        "io"
        "log"
        "net"
        "fmt"
        "errors"
        "crypto/tls"
        "encoding/json"
        "github.com/hlhv/fsock"
        "github.com/hlhv/protocol"
)

type Band struct {
        conn     net.Conn
        reader   *fsock.Reader
        writer   *fsock.Writer
        open     bool
        callback func(*Band, protocol.FrameKind, []byte)
}

func spawnBand (
        address string,
        uuid string,
        callback func(*Band, protocol.FrameKind, []byte),
        tlsConf *tls.Config,
) (
        band *Band,
        err error,
) {
        log.Println("connecting new band")

        log.Println("dialing")
        conn, err := tls.Dial("tcp", address, tlsConf)
        if err != nil { conn.Close(); return nil, err }

        reader := fsock.NewReader(conn)
        writer := fsock.NewWriter(conn)

        log.Println("requesting band status")
        _, err = protocol.WriteMarshalFrame (writer, &protocol.FrameIAm {
                ConnKind: protocol.ConnKindBand,
                Uuid:     uuid,
        })
        if err != nil { conn.Close(); return nil, err }

        kind, data, err := protocol.ReadParseFrame(reader)
        if err != nil { conn.Close(); return nil, err }
        if kind != protocol.FrameKindAccept {
                conn.Close()
                return nil, errors.New (fmt.Sprint (
                        "server sent strange response:", kind))
        }

        frame := protocol.FrameAccept {}
        err = json.Unmarshal(data, &frame)
        if err != nil { conn.Close(); return nil, err }
        log.Println("accepted")

        band = &Band {
                conn:     conn,
                reader:   reader,
                writer:   writer,
                open:     true,
                callback: callback,
        }

        go band.listen()
        return band, nil
}

func (band *Band) listen () {
        for {
                kind, data, err := protocol.ReadParseFrame(band.reader)
                if err == io.EOF { break }
                if err != nil { log.Println("band error:", err) }
                if band.callback == nil {
                        log.Println("BAND CALLBACK NOT REGISTERED")
                } else {
                        band.callback(band, kind, data)
                }
        }
        log.Println("band disconnected")
}

/* Close closes the connection and marks the band as closed so that it can be
 * removed from the list later.
 */
func (band *Band) Close () {
        log.Println("closing band")
        band.open = false
        band.conn.Close()
}

/* ReadParseFrame reads a single frame and parses it, separating the kind and
 * the data.
 */
func (band *Band) ReadParseFrame () (
        kind protocol.FrameKind,
        data []byte,
        err error,
) {
        kind, data, err = protocol.ReadParseFrame(band.reader)
        if err != nil { band.Close() }
        return
}

/* WriteMarshalFrame marshals and writes a Frame.
 */
func (band *Band) WriteMarshalFrame (frame protocol.Frame) (nn int, err error) {
        nn, err = protocol.WriteMarshalFrame(band.writer, frame)
        if err != nil { band.Close() }
        return
}

/* WriteHTTPHead writes HTTP header information. It should only be called once
 * when serving an HTTP response.
 */
func (band *Band) WriteHTTPHead (
        code int,
        headers map[string] []string,
) (
        nn int,
        err error,
) {
        if headers == nil {
                headers = make(map[string] []string)
        }
        return band.WriteMarshalFrame (&protocol.FrameHTTPResHead {
                StatusCode: code,
                Headers:    headers,
        })
}

/* WriteHTTPBody writes a chunk of the response body.
 */
func (band *Band) WriteHTTPBody (data []byte) (nn int, err error) {
        return band.writer.WriteFrame (
                append (
                        []byte { byte(protocol.FrameKindHTTPResBody) },
                        data...
                ),
        )
}

/* writeHTTPEnd ends the HTTP response. This function should be called
 * automatically by the internal callback set by the leash.
 */
func (band *Band) writeHTTPEnd () (nn int, err error) {
        return band.writer.WriteFrame (
                []byte { byte(protocol.FrameKindHTTPResEnd) },
        )
}

/* ReadHTTPBody reads a chunk of the request body. This function returns true
 * for getNext if the chunk was successfully read, and false if it encountered
 * an error or the request ended.
 */
func (band *Band) ReadHTTPBody () (getNext bool, data []byte, err error) {
        getNext = false
        var kind protocol.FrameKind
        kind, data, err = band.ReadParseFrame()
        if err != nil { return }
        switch kind {
        case protocol.FrameKindHTTPReqBody:
                getNext = true
                break
        case protocol.FrameKindHTTPReqEnd:
                data = nil
                break
        default:
                err = errors.New (fmt.Sprint (
                        "got unexpected kind code while processing http req:",
                        kind,
                ))
        }
        return
}
