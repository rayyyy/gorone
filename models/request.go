package models

import (
	"encoding/json"
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
	fmt.Println("未完了のリクエストを取得中")
	requests := GetUnCompletedRequests()
	for _, r := range requests {
		fmt.Printf("RequestID: %v, DesiredCount: %v, NowWorkerNum: %v\n", r.ID, r.DesiredCount, len(r.Assignments))
		if r.DesiredCount > len(r.Assignments) {
			return &r
		}
	}
	// TODO: なかった場合はECSの実行数を適正にしたい
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

// IsDone 終了しているか確認
func (r Request) IsDone() bool {
	db := db.DbManager()
	var count int64
	db.Model(&CalcResult{}).Where("key_name IN ?", r.DecodedBody()).Count(&count)
	return true
}

// DecodedBody bodyをJSONデコードしたものを返す
func (r Request) DecodedBody() []string {
	var data map[string][]string
	err := json.Unmarshal(r.Body, &data)
	if err != nil {
		return []string{}
	}
	return data["values"]
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
