package main

import (
	"io"
	"log"
)

type LogReadWriteCloser struct {
	name   string
	stream io.ReadWriteCloser
}

func newLogReadWriteCloser(name string, stream io.ReadWriteCloser) io.ReadWriteCloser {
	return &LogReadWriteCloser{
		name:   name,
		stream: stream,
	}
}

func (s *LogReadWriteCloser) Close() error {
	log.Println("--->", s.name, "close")
	return s.Close()
}

func (s *LogReadWriteCloser) Read(b []byte) (int, error) {
	log.Println("--->", s.name, "start to read ...")
	l, err := s.stream.Read(b)
	log.Println("--->", s.name, "read len:", l, "data:", b[:l])
	return l, err
}

func (s *LogReadWriteCloser) Write(b []byte) (int, error) {
	log.Println("--->", s.name, "start to write len:", len(b), "data:", b)
	l, err := s.stream.Write(b)
	log.Println("--->", s.name, "write err", err)
	return l, err
}
