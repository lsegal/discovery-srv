package registry

import (
	"errors"
	"sync"

	proto "github.com/micro/go-platform/discovery/proto"
)

type watcher struct {
	once  sync.Once
	watch chan *proto.Result
	stop  chan bool
}

var (
	ErrWatchStopped = errors.New("watch stopped")
)

func newWatcher(watch chan *proto.Result, stop chan bool) *watcher {
	var once sync.Once
	return &watcher{
		once:  once,
		watch: watch,
		stop:  stop,
	}
}

func (w *watcher) Next() (*proto.Result, error) {
	r, ok := <-w.watch
	if !ok {
		return nil, ErrWatchStopped
	}
	return r, nil
}

func (w *watcher) Stop() {
	w.once.Do(func() {
		close(w.stop)
	})
}
