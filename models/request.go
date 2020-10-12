package models

import (
	"fmt"
	"gorone/db"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Request タスク
type Request struct {
	RedisTagName string
	JobType      string
	DesiredCount int
	gorm.Model
	Body        datatypes.JSON
	DoneAt      *time.Time
	Assignments []Assignment
}

// FindRequest workerが足りないRequestを一つ取得
func FindRequest() *Request {
	requests := GetUnCompletedRequests()
	for _, r := range requests {
		fmt.Printf("id:%s c:%s d:%s\n", r.ID, r.DesiredCount, len(r.Assignments))
		r.Assignments[0].Terminate()
		if r.DesiredCount > len(r.Assignments) {
			// r.Done()
			return &r
		}
	}
	// なかった場合はECSの実行数を適正にしたい
	return nil
}

// GetUnCompletedRequests 未完了のリクエストを取得
func GetUnCompletedRequests() []Request {
	db := db.DbManager()
	var requests []Request
	if err := db.Preload("Assignments", "terminated_at IS NULL").Where("done_at IS NULL").Find(&requests).Error; err != nil {
		return nil
	}
	return requests
}

// Done Requestが完了したとき
func (r Request) Done() error {
	db := db.DbManager()
	now := time.Now()
	r.DoneAt = &now

	if err := db.Save(&r).Error; err != nil {
		return err
	}

	return nil
}
