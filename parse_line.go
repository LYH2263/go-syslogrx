package syslogrx

import "fmt"

func ParseLine(line string) (*Message, error) {
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}
	if m, err := ParseRFC5424(line); err == nil {
		return m, nil
	}
	m, err := ParseRFC3164(line)
	if err != nil {
		return nil, err
	}
	return m, nil
}
