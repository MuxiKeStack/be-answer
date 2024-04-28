package domain

import "time"

type Answer struct {
	Id          int64
	PublisherId int64
	QuestionId  int64
	Content     string
	Utime       time.Time
	Ctime       time.Time
}
