package models

import "time"

type Article struct {
    ID          int       `json:"id"`
    Title       string    `json:"title" binding:"required,min=20"`
    Content     string    `json:"content" binding:"required,min=200"`
    Category    string    `json:"category" binding:"required,min=3"`
    CreatedDate time.Time `json:"created_date"`
    UpdatedDate time.Time `json:"updated_date"`
    Status      string    `json:"status" binding:"required,oneof=publish draft trash"`
}