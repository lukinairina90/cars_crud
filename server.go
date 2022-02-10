package main

import (
	"C_CRUD/models"
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type ViewData struct {
	Title   string
	Message string
}

type crud struct {
	db  *gorm.DB
	cfg Config
}

func main() {
	// https://github.com/go-ozzo/ozzo-validation
	// form method = POST , action = /api/cars
	// form method = POST , action = /api/cars/{id}
	//https://github.com/go-chi/chi
	//https://www.newline.co/@kchan/building-a-simple-restful-api-with-go-and-chi--5912c411
	// /api/cars/crudmethhods
	// GET /api/cars index
	// POST /api/cars store
	// POST /api/cars/{id} update
	// DELETE /api/cars/{id} delete (form method delete)
	// /web/cars/index

	// describe routes
	// add gorm, describe car model, run migrations
	// start POST /api/cars method

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	data := ViewData{
	//		Title:   "Manage cars",
	//		Message: "Cars",
	//	}
	//	tmpl, _ := template.ParseFiles("templates/index.html")
	//	err := tmpl.Execute(w, data)
	//	if err != nil {
	//		return
	//	}
	//})

	//ENV
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)

	//GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	db.Logger = newLogger

	cr := crud{
		db:  db,
		cfg: cfg,
	}

	r := chi.NewRouter()
	// RESTy routes for "articles" resource
	r.Route("/api", func(r chi.Router) {
		r.Route("/cars", func(r chi.Router) {
			r.Get("/", cr.listCars)          // GET /api/cars index
			r.Put("/{id}", cr.updateCarById) // UPDATE /api/cars index
			r.Post("/{id}", cr.listCarByID)  // POST /api/cars/{id} get
			r.Post("/", cr.createCar)        // POST /api/cars/ create
			r.Delete("/{id}", cr.deleteCar)  // /api/cars/{id} delete (form method delete)
		})
	})

	// RESTy routes for "articles" resource
	r.Route("/cars", func(r chi.Router) {
		r.Get("/", webIndex) // GET /api/cars index
	})

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	data := ViewData{
	//		Title:   "Manage cars",
	//		Message: "Cars",
	//	}
	//	tmpl, _ := template.ParseFiles("templates/index.html")
	//	err := tmpl.Execute(w, data)
	//	if err != nil {
	//		return
	//	}
	//})
	//
	//r := chi.NewRouter()
	//r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//	_, err := w.Write([]byte("welcome"))
	//	if err != nil {
	//		return
	//	}
	//})
	//// RESTy routes for "articles" resource
	//r.Route("/articles", func(r chi.Router) {
	//	r.Get("/", listCars) // GET /api/cars index
	//	//r.Post("/{id}", listCarByID) // POST /api/cars/{id} edit
	//	//r.Post("/", createCar)       // POST /api/cars/ create
	//	//r.Delete("/", deleteCar)     // /api/cars/{id} delete (form method delete)
	//
	//})

	//Migration
	if err = db.AutoMigrate(&models.Car{}); err != nil {
		return
	}

	err = http.ListenAndServe(":8181", r)
	if err != nil {
		return
	}
}

func (cr crud) deleteCar(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "id")

	logrus.WithField("car_id", key).Info("starting deletion car")

	var car models.Car
	//id, _ := strconv.ParseInt(key, 10, 64)
	if err := cr.db.Debug().Where("id = ?", key).Delete(&car).Error; err != nil {
		logrus.WithField("error", err).Error("error deleting car")
	}

	//cr.db.Delete(&car)
	w.WriteHeader(http.StatusNoContent)
}

func (cr crud) listCarByID(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "id")

	var car models.Car
	cr.db.First(&car, key)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(car)
	if err != nil {
		return
	}
}

func (cr crud) createCar(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var car models.Car
	err := json.Unmarshal(requestBody, &car)
	if err != nil {
		return
	}

	cr.db.Create(&car)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&car)
	if err != nil {
		return
	}
}

func webIndex(w http.ResponseWriter, r *http.Request) {
	data := ViewData{
		Title:   "Manage cars",
		Message: "Cars",
	}
	tmpl, _ := template.ParseFiles("templates/index.html")
	err := tmpl.Execute(w, data)
	if err != nil {
		return
	}
}

func (cr crud) listCars(w http.ResponseWriter, r *http.Request) {
	var carsModels []models.Car
	cr.db.Find(&carsModels)

	logrus.WithField("ListCar", carsModels).Info("starting ListCar")

	carsResp := make([]Car, 0, len(carsModels))
	for _, carModel := range carsModels {
		carsResp = append(carsResp, Car{
			ID:        carModel.ID,
			ModelName: carModel.ModelName,
			Type:      carModel.Type,
			//ModelInfo:    fmt.Sprintf("%s (%s)", carModel.ModelName, carModel.Type),
			Transmission: carModel.Transmission,
			Engine:       carModel.Engine,
			HorsePower:   carModel.HorsePower,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(carsResp)
	if err != nil {
		return
	}
}

func (cr crud) updateCarById(w http.ResponseWriter, r *http.Request) {
	//key := chi.URLParam(r, "id")

	requestBody, _ := ioutil.ReadAll(r.Body)
	var car models.Car

	logrus.WithField("updateCarById", car).Info("starting updateCarById")

	err := json.Unmarshal(requestBody, &car)
	if err != nil {
		return
	}
	cr.db.Save(&car)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(car)
	if err != nil {
		return
	}
}

func createCars(writer http.ResponseWriter, request *http.Request) {
	//request.Form.Get("sefe")
}
