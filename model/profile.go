package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model

	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Nickname  string `json:"nickname"`
	Active    bool   `json:"active"`

	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`

	Faces []Face `json:"faces" gorm:"foreignKey:ProfileId"`
}
