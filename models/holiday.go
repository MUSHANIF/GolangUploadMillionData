package models

import "time"
//harus huruf besar untuk field struct agar bisa diakses oleh gorm
//untuk menghindari error gorm: unsupported driver type []uint8
type Holiday struct {
    Id              uint      `json:"id" gorm:"primaryKey"`
    holiday_name     string    `json:"holiday_name" gorm:"unique;not null"`
    date      string    `json:"date" gorm:"unique;not null"`
    created_at       time.Time `json:"created_at"`
    updated_at       time.Time `json:"updated_at"`
}
