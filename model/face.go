package model

import "gorm.io/gorm"

type Face struct {
	gorm.Model

	ProfileId    uint      `json:"profileId"`
	Profile      *Profile  `json:"profile" gorm:"foreignKey:ProfileId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FaceEncoding []float64 `json:"faceEncoding" gorm:"type:json"`
}
