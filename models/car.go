package models

import "gorm.io/gorm"

type Car struct {
	gorm.Model
	ModelName    string
	Type         string
	Transmission string
	Engine       string
	HorsePower   string
}
