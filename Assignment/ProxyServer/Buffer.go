package main

import (
	"container/list"
	"fmt"
)

const (
	B  int = 1
	KB     = 1024 * B
	MB     = 1024 * KB
	GB     = 1024 * MB
)

var FileBuffer buffer = initBuffer(2 * MB)

type value struct {
	url  string
	file *[]byte
}

type buffer struct {
	list     list.List
	maxSize  int
	freeSize int
}

func initBuffer(size int) buffer {
	if size < 0 {
		size = 0
	}
	buf := buffer{
		list:     *list.New().Init(),
		maxSize:  size,
		freeSize: size,
	}
	return buf
}

func (buf *buffer) addFile(url string, file []byte) error {
	cap := cap(file)
	if cap > buf.maxSize {
		return fmt.Errorf("file capacity %d bigger than buf max size %d", cap, buf.maxSize)
	}
	for cap > buf.freeSize {
		buf.removeFile()
	}
	val := value{
		url:  url,
		file: &file,
	}
	buf.list.PushFront(val)
	buf.freeSize -= cap
	return nil
}

func (buf *buffer) searchAndUpdateBuffer(url string) []byte {
	node := buf.list.Front()
	for i := 0; i < buf.list.Len(); i++ {
		if node.Value.(value).url == url {
			buf.list.MoveToFront(node)
			return *node.Value.(value).file
		}
		node = node.Next()
	}
	return nil
}

func (buf *buffer) removeFile() error {
	if buf.freeSize == buf.maxSize {
		return fmt.Errorf("buffer empty, free size %d", buf.freeSize)
	}
	buf.list.Remove(buf.list.Back())
	return nil
}
