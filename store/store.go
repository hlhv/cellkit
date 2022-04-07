package store

import (
        "errors"
        "io/ioutil"
        "github.com/hlhv/protocol"
        "github.com/hlhv/cellkit/client"
)

/* Store is a simple resource manager for serving static file resources. Files
 * can be registered and unregistered dynamically, and are loaded lazily. It can
 * be combined with any other system for serving files and pages.
 */
type Store struct {
        items map[string] *LazyFile
        root  string
}

/* New creates a new Store.
 */
func New (root string) (store *Store) {
        lastIndex := len(root) - 1
        if root[lastIndex] == '/' {
                root = root[:lastIndex]
        }
        return &Store {
                items: make(map[string] *LazyFile),
                root:  root,
        }
}

/* Register registers a file located at the filepath on the specific url path.
 * Url paths must start from /, and not end in /.
 */
func (store *Store) Register (filePath string, webPath string) (err error) {
        if filePath[0] != '/' {
                return errors.New("file path must start at /")
        }
        
        filePath = store.root + filePath
        
        if webPath[0] != '/' {
                return errors.New("web path must start at /")
        }
        
        if webPath[len(webPath) - 1] == '/' {
                return errors.New("web path must be a file, not a directory")
        }
        
        store.items[webPath] = &LazyFile { FilePath: filePath }
        return nil
}

/* RegisterDir registers every file in a given directory.
 */
func (store *Store) RegisterDir (dirPath string, webPath string) (err error) {
        if dirPath[0] != '/' {
                return errors.New("dir path must start at /")
        }
        
        dirPath = store.root + dirPath
        
        if webPath[0] != '/' {
                return errors.New("web path must start at /")
        }
        
        if dirPath[len(dirPath) - 1] != '/' { dirPath += "/" }
        if webPath[len(webPath) - 1] != '/' { webPath += "/" }
        
        directory, err := ioutil.ReadDir(dirPath)
        if err != nil { return err }

        for _, file := range(directory) {
                if file.IsDir() { continue }
                store.items[webPath + file.Name()] = &LazyFile {
                        FilePath: dirPath + file.Name(),
                }
        }
        return nil
}

/* Unregister finds the file registered at the specified url path and
 * unregisters it, freeing it from memory
 */
func (store *Store) Unregister (webPath string) (err error) {
        _, exists := store.items[webPath]
        if !exists {
                return errors.New("path " + webPath + " is not registered")
        }
        delete(store.items, webPath)
        return nil
}

/* TryHandle checks the request path against the map of registered files, and
 * serves a match if it finds it. The function returns wether it served a file
 * or not. If this function returns false, the request needs to be handled
 * still.
 */
func (store *Store) TryHandle (
        band *client.Band,
        head *protocol.FrameHTTPReqHead,
) (
        handled bool,
        err     error,
) {
        item, matched := store.items[head.Path]
        if !matched { return false, nil }
        err = item.Send(band, head)
        return true, err
}

/* Returns the root path of the store. This can be helpful for doing things such
 * as registering an entire directory while doing operations on the files inside
 * of it.
 */
func (store *Store) GetRoot () (root string) {
        return store.root
}
