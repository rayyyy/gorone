package models

import (
	"gorm.io/gorm"
)

// CalcResult タスク
type CalcResult struct {
	gorm.Model
	KeyName string
	Result  int
}
