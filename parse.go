package syslogrx

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Priority  int
	Facility  int
	Severity  int
	Timestamp time.Time
	Hostname  string
	Tag       string
	Content   string
	Raw       string
	RawBytes  []byte
}

func ParseRFC3164(line string) (*Message, error) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "<") {
		return nil, ErrInvalid
	}
	end := strings.IndexByte(line, '>')
	if end < 0 {
		return nil, ErrInvalid
	}
	pri, err := strconv.Atoi(line[1:end])
	if err != nil {
		return nil, err
	}
	rest := strings.TrimSpace(line[end+1:])
	if len(rest) < 15 {
		return nil, fmt.Errorf("%w: short", ErrInvalid)
	}
	tsStr := rest[:15]
	ts, err := time.ParseInLocation("Jan  2 15:04:05", tsStr, time.UTC)
	if err != nil {
		ts, err = time.ParseInLocation("Jan 2 15:04:05", strings.Join(strings.Fields(tsStr), " "), time.UTC)
		if err != nil {
			return nil, err
		}
	}
	rest = strings.TrimSpace(rest[15:])
	parts := strings.SplitN(rest, " ", 3)
	if len(parts) < 2 {
		return nil, ErrInvalid
	}
	host, tag := parts[0], parts[1]
	content := ""
	if len(parts) == 3 {
		content = parts[2]
	}
	tag = strings.TrimSuffix(tag, ":")
	return &Message{
		Priority:  pri,
		Facility:  pri / 8,
		Severity:  pri % 8,
		Timestamp: ts,
		Hostname:  host,
		Tag:       tag,
		Content:   content,
		Raw:       line,
	}, nil
}
