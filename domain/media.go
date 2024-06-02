package domain

import "time"

type Media struct {
	MediaId   int64
	UserId    int64
	FilePath  string
	FileType  string
	CreatedAt time.Time
}
