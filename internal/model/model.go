package model

import (
	"time"
)
//one line url string
type Requestline struct {
	Nline    int64
	Url      string
	Readtime time.Duration
	Size     int64
	Err      error
}