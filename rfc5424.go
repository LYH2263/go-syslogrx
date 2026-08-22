package syslogrx
import (
        "fmt"
        "strconv"
        "strings"
        "time"
)
func ParseRFC5424(line string) (*Message, error) {
        line = strings.TrimSpace(line)
        if !strings.HasPrefix(line, "<") { return nil, ErrInvalid }
        end := strings.IndexByte(line, '>')
        if end < 0 { return nil, ErrInvalid }
        pri, err := strconv.Atoi(line[1:end])
        if err != nil { return nil, err }
        rest := line[end+1:]
        // VERSION SP TIMESTAMP SP HOSTNAME SP APP-NAME SP PROCID SP MSGID SP STRUCTURED-DATA [SP MSG]
        parts := strings.SplitN(rest, " ", 7)
        if len(parts) < 6 { return nil, fmt.Errorf("%w: fields", ErrInvalid) }
        ts := time.Time{}
        if parts[1] != "-" {
                ts, err = time.Parse(time.RFC3339Nano, parts[1])
                if err != nil {
                        ts, err = time.Parse(time.RFC3339, parts[1])
                        if err != nil { return nil, err }
                }
        }
        host, app := parts[2], parts[3]
        content := ""
        if len(parts) >= 7 { content = parts[6] }
        return &Message{
                Priority: pri, Facility: pri / 8, Severity: pri % 8,
                Timestamp: ts, Hostname: host, Tag: app, Content: content, Raw: line,
        }, nil
}
