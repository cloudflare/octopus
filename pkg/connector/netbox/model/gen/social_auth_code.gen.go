// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameSocialAuthCode = "social_auth_code"

// SocialAuthCode mapped from table <social_auth_code>
type SocialAuthCode struct {
	ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Email     string    `gorm:"column:email;not null" json:"email"`
	Code      string    `gorm:"column:code;not null" json:"code"`
	Verified  bool      `gorm:"column:verified;not null" json:"verified"`
	Timestamp time.Time `gorm:"column:timestamp;not null" json:"timestamp"`
}

// TableName SocialAuthCode's table name
func (*SocialAuthCode) TableName() string {
	return TableNameSocialAuthCode
}
