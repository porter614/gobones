package watch

import (
	"errors"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
)

// WatchFile starts a file watcher on the specified path and returns a channel
// for communicating config changes. When the watched path changes, the config
// is reloaded and pushed on to the channel.
func WatchFile(filepath string) (chan []byte, chan error, error) {
	ch := make(chan []byte)
	errch := make(chan error)

	// Create a filesystem watcher to monitor the config file
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		if watcher != nil {
			watcher.Close()
		}
		return ch, errch, err
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case e, ok := <-watcher.Events:
				if !ok {
					errch <- errors.New("Watcher not OK")
					if err := watcher.Remove(filepath); err != nil {
						watcher.Close()
						return
					}
					if err := watcher.Add(filepath); err != nil {
						watcher.Close()
						return
					}
					continue
				}
				if (e.Op&fsnotify.Write == fsnotify.Write ||
					e.Op&fsnotify.Create == fsnotify.Create ||
					e.Op&fsnotify.Remove == fsnotify.Remove) &&
					e.Name == filepath {
					// Read the file and put the contents on the channel
					content, err := ioutil.ReadFile(filepath)
					if err != nil {
						errch <- err
						if err := watcher.Remove(filepath); err != nil {
							watcher.Close()
							return
						}
						if err := watcher.Add(filepath); err != nil {
							watcher.Close()
							return
						}
						continue
					}
					// Need to re-watch the file for kubernetes configmaps
					// since the symlink was updated
					if err := watcher.Remove(filepath); err != nil {
						watcher.Close()
						return
					}
					if err := watcher.Add(filepath); err != nil {
						watcher.Close()
						return
					}
					ch <- content
				}
			case err := <-watcher.Errors:
				errch <- err
				if err := watcher.Remove(filepath); err != nil {
					watcher.Close()
					return
				}
				if err := watcher.Add(filepath); err != nil {
					watcher.Close()
					return
				}
				continue
			}
		}
	}()

	// Tell the watcher to watch
	if err := watcher.Add(filepath); err != nil {
		watcher.Close()
		return ch, errch, err
	}
	return ch, errch, nil
}
