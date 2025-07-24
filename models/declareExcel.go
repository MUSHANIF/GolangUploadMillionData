package models

import "time"
//harus huruf besar untuk field struct agar bisa diakses oleh gorm
//untuk menghindari error gorm: unsupported driver type []uint8
type DeclareExcel struct {
    Id              uint      `json:"id" gorm:"primaryKey"`
    ConoteDate      string    `json:"conote_date"`
    ProductCode     string    `json:"product_code" gorm:"unique;not null"`
    NoStikb         string    `json:"no_stikb" gorm:"unique;not null"`
    CustName        string    `json:"cust_name"`
    CustNo          string    `json:"cust_no"`
    NoResi          string    `json:"no_resi"`
    Origin          string    `json:"origin"`
    Destination     string    `json:"destination"`
    Description     string    `json:"description"`
    Status      int `gorm:column:status;`
	Rate			string   `json:"rate"`
    SumInsured      string    `json:"sum_insured"`
    Premium         string    `json:"premium"`
    ChangeStatusFrom int `gorm:"default:0"` 
    Reason          string    `json:"reason"`
    UpdatedBy *int `gorm:"column:updated_by;default:0"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
