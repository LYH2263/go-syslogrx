package syslogrx

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileSink struct {
	mu   sync.Mutex
	path string
	f    *os.File
}

func OpenFileSink(path string) (*FileSink, error) {
	if d := filepath.Dir(path); d != "." {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &FileSink{path: path, f: f}, nil
}

func (s *FileSink) Write(m *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.f == nil {
		return fmt.Errorf("filesink closed")
	}
	_, err := s.f.Write(append(append([]byte(nil), m.RawBytes...), 10))
	return err
}

func (s *FileSink) Flush() []*Message { return nil }
func (s *FileSink) Clear()            {}

func (s *FileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.f == nil {
		return nil
	}
	err := s.f.Close()
	s.f = nil
	return err
}

func (s *FileSink) Rotate(newPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d := filepath.Dir(newPath); d != "." {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if s.f != nil {
		_ = s.f.Close()
	}
	s.f = f
	s.path = newPath
	return nil
}
