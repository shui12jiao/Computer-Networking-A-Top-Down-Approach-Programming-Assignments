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

var Buffer buffer

type value struct {
	url  string
	file *[]byte
}

type buffer struct {
	list     list.List
	maxSize  int
	freeSize int
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

//TODO
func (buf *buffer) updateBuffer(url string, file []byte) error {
	e := list.Element{
		Value: value{
			url:  url,
			file: &file,
		},
	}
	buf.list.MoveToFront(&e)
	if buf.list.Front() != &e {
		return fmt.Errorf("update fail")
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

//TODO
