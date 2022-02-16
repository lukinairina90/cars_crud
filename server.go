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
	"html/template"
	"io/ioutil"
	"net/http"
)

type ViewData struct {
	gorm.Model
	ID           int
	ModelName    string
	Type         string
	Transmission string
	Engine       string
	HorsePower   string
}

type crud struct {
	db  *gorm.DB
	cfg Config
}

func main() {
	// https://github.com/go-ozzo/ozzo-validation

	//ENV
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold:             time.Second,   // Slow SQL threshold
	//		LogLevel:                  logger.Silent, // Log level
	//		IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
	//		Colorful:                  true,          // Disable color
	//	},
	//)

	//GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	//db.Logger = newLogger

	cr := crud{
		db:  db,
		cfg: cfg,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/cars", func(r chi.Router) {
			r.Get("/", cr.listCars)
			r.Post("/", cr.createCar)
			r.Put("/{id:[0-9]+}", cr.updateCarById)
			r.Get("/{id:[0-9]+}", cr.listCarByID)
			r.Delete("/{id:[0-9]+}", cr.deleteCar)
		})
	})
	r.Route("/cars", func(r chi.Router) {
		r.Get("/", cr.webListCars)
	})

	//Migration
	if err = db.AutoMigrate(&models.Car{}); err != nil {
		return
	}

	err = http.ListenAndServe(":8181", r)
	if err != nil {
		return
	}
}

func (cr crud) getCarsList() []Car {
	var carsModels []models.Car
	cr.db.Find(&carsModels)

	logrus.WithField("ListCar", carsModels).Info("starting ListCar")

	carsResp := make([]Car, 0, len(carsModels))
	for _, carModel := range carsModels {
		carsResp = append(carsResp, Car{
			ID:           carModel.ID,
			ModelName:    carModel.ModelName,
			Type:         carModel.Type,
			Transmission: carModel.Transmission,
			Engine:       carModel.Engine,
			HorsePower:   carModel.HorsePower,
			//ModelInfo:    fmt.Sprintf("%s (%s)", carModel.ModelName, carModel.Type),
		})
	}
	return carsResp
}

func (cr crud) webListCars(w http.ResponseWriter, r *http.Request) {
	carsResp := cr.getCarsList()
	tmpl, _ := template.ParseFiles("templates/index.html")
	err := tmpl.Execute(w, carsResp)
	if err != nil {
		println(err)
	}
}

func (cr crud) listCars(w http.ResponseWriter, r *http.Request) {
	carsResp := cr.getCarsList()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(carsResp)
	if err != nil {
		return
	}
}

type CreateCarRequest struct {
	ModelType    string `json:"model_type"`
	Type         string `json:"type"`
	Transmission string `json:"transmission"`
	Engine       string `json:"engine"`
	HorsePower   string `json:"horse_power"`
}

func (cr crud) createCar(response http.ResponseWriter, request *http.Request) {

	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := CreateCarRequest{}
	if err := json.Unmarshal(b, &req); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	car := models.Car{
		ModelName:    req.ModelType,
		Type:         req.Type,
		Transmission: req.Transmission,
		Engine:       req.Engine,
		HorsePower:   req.HorsePower,
	}

	if err := cr.db.Create(&car).Error; err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(response).Encode(&car)
	if err != nil {
		return
	}
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

func (cr crud) updateCarById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	requestBody, _ := ioutil.ReadAll(r.Body)
	var ca CarUpdate

	logrus.WithField("updateCarById", ca).Info("starting updateCarById")

	err := json.Unmarshal(requestBody, &ca)
	if err != nil {
		return
	}

	model := models.Car{
		ModelName:    ca.ModelName,
		Type:         ca.Type,
		Transmission: ca.Transmission,
		Engine:       ca.Engine,
		HorsePower:   ca.HorsePower,
	}

	if err := cr.db.Where("id = ?", id).Updates(model); err == nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ca)
	if err != nil {
		return
	}
}

func (cr crud) deleteCar(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "id")

	logrus.WithField("car_id", key).Info("starting deletion car")

	var car models.Car
	if err := cr.db.Debug().Where("id = ?", key).Delete(&car).Error; err != nil {
		logrus.WithField("error", err).Error("error deleting car")
	}

	w.WriteHeader(http.StatusNoContent)
}
