package model

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;"`
	Email    string `json:"email" validate:"required,email" gorm:"unique;not null"`
	Password string `json:"password" validate:"required"`
	Remember string `json:"remember"`
	
}

type Post struct {
	ID 			uint `json:"id" gorm:"primaryKey;"`
	UserID 		uint 	
	User		User	`gorm:"foreignKey:UserID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title 		string    `json:"title" gorm:"not null;"`
	Content 	string  `json:"content" gorm:"not null;"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
 }

type Like struct{
	ID		uint `json:"id"`
	PostID 	uint `json:"postid" gorm:"index:idx_name,unique"`
	UserID 	uint `json:"userid" gorm:"index:idx_name,unique"`
	User	User `gorm:"foreignKey:UserID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post	Post `gorm:"foreignKey:PostID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CntList struct {
	Pid		uint	`json:"pid" gorm:"postid;"`
	Cnt		int64	`json:"cnt" gorm:"cnt;"`
}

 func UsernewToken() string{
    b := make([]byte, 10)
    rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
 }

