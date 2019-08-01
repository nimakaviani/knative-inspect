package tests

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"sync"
)

type buffer struct {
	contents   []byte
	readCursor uint64
	lock       *sync.Mutex
	closed     bool
}

func BufferReader(reader io.Reader) *buffer {
	b := &buffer{
		lock: &sync.Mutex{},
	}

	io.Copy(b, reader)
	b.Close()

	return b
}

func (b *buffer) Write(p []byte) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.closed {
		return 0, fmt.Errorf("attempt to write to closed buffer")
	}

	println("here!!")

	b.contents = append(b.contents, p...)
	return len(p), nil
}

func (b *buffer) Read(d []byte) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.closed {
		return 0, fmt.Errorf("attempt to read from closed buffer")
	}

	if uint64(len(b.contents)) <= b.readCursor {
		return 0, io.EOF
	}

	n := copy(d, b.contents[b.readCursor:])
	b.readCursor += uint64(n)

	return n, nil
}

func (b *buffer) Close() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.closed = true

	return nil
}

func (b *buffer) Closed() bool {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.closed
}

func (b *buffer) Contents() []byte {
	b.lock.Lock()
	defer b.lock.Unlock()

	contents := make([]byte, len(b.contents))
	copy(contents, b.contents)
	return contents
}

func (b *buffer) ShouldSay(expr string) (bool, []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()

	re := regexp.MustCompile(expr)

	unreadBytes := b.contents[b.readCursor:]
	copyOfUnreadBytes := make([]byte, len(unreadBytes))
	copy(copyOfUnreadBytes, unreadBytes)

	loc := re.FindIndex(unreadBytes)

	if loc != nil {
		b.readCursor += uint64(loc[1])
		return true, copyOfUnreadBytes
	}
	log.Fatalf("didnot find %s", expr)
	return false, copyOfUnreadBytes
}
