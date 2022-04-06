package store

import (
        "os"
        "io"
        "strings"
        "net/http"
        "path/filepath"
        "github.com/hlhv/protocol"
        "github.com/hlhv/cellkit/client"
)

/* chunkSize does not refer to actual chunked encoding. This is just so the
 * client doesn't have to wait for the cell to send everything over and the
 * queen to send that over before recieving everything. It should be at least
 * 512 in order for accurate mime-type detection.
 */
const chunkSize int = 1024

/* LazyFile is a struct capable of serving a file. The file is cached into
 * memory when it is first loaded, hence the name.
 */
type LazyFile struct {
        FilePath string
        mime     string
        chunks   []fileChunk
}

type fileChunk []byte

/* Send sends the file along with a content-type header.
 */
func (item *LazyFile) Send (
        band *client.Band,
        head *protocol.FrameHTTPReqHead,
) (
        err error,
) {
        if item.chunks == nil {
                err = item.loadAndSend(band, head)
                return err
        }
        
        _, err = band.WriteHTTPHead(200, map[string] []string{
                "content-type": []string { item.mime },
        })
        if err != nil { return err }

        for _, chunk := range(item.chunks) {
                _, err = band.WriteHTTPBody(chunk)
                if err != nil { return err }
        }

        return nil
}

/* loadAndSend loads the file from disk while sending it in response to an http
 * request. This should be called when there is an http request for this file
 * but it has not been loaded yet.
 */
func (item *LazyFile) loadAndSend (
        band *client.Band,
        head *protocol.FrameHTTPReqHead,
) (
        err error,
) {
        file, err := os.Open(item.FilePath)
        defer file.Close()
        if err != nil { return err }
        

        needMime := true
        for {
                chunk := make([]byte, chunkSize)
                bytesRead, err := io.ReadFull(file, chunk)
                chunk = chunk[:bytesRead]

                fileEnded := err == io.ErrUnexpectedEOF || err == io.EOF
		if err != nil && !fileEnded {
                        return err
                }

                if needMime {
                        needMime = false
                        item.mime = mimeSniff(item.FilePath, chunk)
                        _, err = band.WriteHTTPHead(200, map[string] []string{
                                "content-type": []string { item.mime },
                        })
                        if err != nil { return err }
                }

                item.chunks = append(item.chunks, chunk)
                band.WriteHTTPBody(chunk)
		
                if fileEnded { break }
        }
        
        return nil
}

/* mimeSniff determines the content type of a byte array and an associated name.
 * This isn't very good as of now but it works!
 */
func mimeSniff (name string, data []byte) (mime string) {
        extension := filepath.Ext(name)
        mime = http.DetectContentType(data)
        if (
                strings.HasPrefix(mime, "text/plain") &&
                extension != ".txt" &&
                extension != "" ) {
                
                mime = strings.Replace(mime, "plain", extension[1:], 1)
        }
        return mime
}
