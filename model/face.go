package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Face struct {
	gorm.Model

	ProfileId    uint            `json:"profileId"`
	Profile      *Profile        `json:"profile" gorm:"foreignKey:ProfileId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FaceEncoding pq.Float64Array `json:"faceEncoding" gorm:"type:float8[]"` // `float8[]` cho PostgreSQL array
}
