package syslogrx
import "errors"
var (
        ErrClosed = errors.New("syslogrx: closed")
        ErrInvalid = errors.New("syslogrx: invalid")
        ErrNotFound = errors.New("syslogrx: not found")
        ErrConflict = errors.New("syslogrx: conflict")
)
