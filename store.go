package syslogrx
import (
        "bufio"
        "os"
        "sync"
)
type Store struct {
        mu sync.Mutex
        path string
        f *os.File
        w *bufio.Writer
        msgs []*Message
        closed bool
}
func OpenStore(path string) (*Store, error) {
        f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
        if err != nil { return nil, err }
        return &Store{path: path, f: f, w: bufio.NewWriter(f)}, nil
}
func (s *Store) Append(m *Message) error {
        s.mu.Lock(); defer s.mu.Unlock()
        if s.closed { return ErrClosed }
        if s.w == nil { return ErrInvalid }
        if _, err := s.w.WriteString(m.Raw + "\n"); err != nil { return err }
        cp := *m
        s.msgs = append(s.msgs, &cp)
        return nil
}
func (s *Store) Flush() error {
        s.mu.Lock(); defer s.mu.Unlock()
        if s.w == nil { return nil }
        if err := s.w.Flush(); err != nil { return err }
        return s.f.Sync()
}
func (s *Store) Recent(n int) []*Message {
        s.mu.Lock(); defer s.mu.Unlock()
        if n <= 0 || n > len(s.msgs) { n = len(s.msgs) }
        out := make([]*Message, n)
        copy(out, s.msgs[len(s.msgs)-n:])
        return out
}
func (s *Store) Close() error {
        s.mu.Lock(); defer s.mu.Unlock()
        if s.closed { return nil }
        s.closed = true
        if s.w != nil { _ = s.w.Flush() }
        if s.f != nil {
                err := s.f.Close()
                s.f = nil
                s.w = nil
                return err
        }
        return nil
}
