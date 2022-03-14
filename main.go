package main

import (
	"C_CRUD/car"
	"C_CRUD/configuration"
	"C_CRUD/models"
	"C_CRUD/routes"
	"fmt"
	"github.com/caarlos0/env/v6"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	// Parsing env's.
	cfg := configuration.Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	// Creating GORM connection.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	carTransport := car.NewTransport(db, cfg)

	router := routes.InitRouter(carTransport)

	// Run GORM auto migrations.
	if err = db.AutoMigrate(&models.Car{}); err != nil {
		return
	}

	err = http.ListenAndServe(":8181", router)
	if err != nil {
		return
	}
}
