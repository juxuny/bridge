package log

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

type writer struct {
	dir      string
	prefix   string
	f        *os.File
	buf      *bufio.Writer
	lastDate string
	lock     *sync.Mutex
}

func NewWriter(dir string, prefix string) *writer {
	w := &writer{
		prefix:   prefix,
		dir:      dir,
		lastDate: time.Now().Format(dateFormat),
		lock:     &sync.Mutex{},
	}
	w.initBufWriter()
	return w
}

func (t *writer) Flush() error {
	return t.buf.Flush()
}

func (t *writer) initBufWriter() {
	var err error
	t.f, err = os.OpenFile(t.dir+string(os.PathSeparator)+fmt.Sprintf("%s-%s.log", t.prefix, t.lastDate), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		DefaultLogger.Error(err)
	}
	t.buf = bufio.NewWriter(t.f)
}

func (t *writer) Write(p []byte) (n int, err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	currentDate := time.Now().Format(dateFormat)
	if currentDate != t.lastDate {
		t.lastDate = currentDate
		if err := t.buf.Flush(); err != nil {
			DefaultLogger.Error(err)
		}
		_ = t.f.Close()
		t.initBufWriter()
	}
	return t.buf.Write(p)
}
