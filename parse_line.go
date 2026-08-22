package syslogrx

func ParseLine(line string) (*Message, error) {
	if line == "" {
		return nil, Malformed(ErrInvalid)
	}
	if m, err := ParseRFC5424(line); err == nil {
		return m, nil
	}
	m, err := ParseRFC3164(line)
	if err != nil {
		return nil, Malformed(err)
	}
	return m, nil
}
