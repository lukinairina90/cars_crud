package car

import (
	"C_CRUD/car/repositories"
	"C_CRUD/configuration"
	"C_CRUD/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
)

func NewTransport(carRepository *repositories.CarRepository, cfg configuration.Config) *Transport {
	return &Transport{
		carRepository: carRepository,
		cfg:           cfg,
	}
}

type Transport struct {
	carRepository *repositories.CarRepository
	cfg           configuration.Config
}

func (t Transport) ShowCars(w http.ResponseWriter, r *http.Request) {
	carModels, err := t.carRepository.List()
	if err != nil {
		FireInternalServerError(w, "error getting cars list", err)
	}

	cars := make([]Car, 0, len(carModels))
	for _, carModel := range carModels {
		cars = append(cars, Car{
			ID:           carModel.ID,
			ModelName:    carModel.ModelName,
			Type:         carModel.Type,
			Transmission: carModel.Transmission,
			Engine:       carModel.Engine,
			HorsePower:   carModel.HorsePower,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cars); err != nil {
		FireInternalServerError(w, "error encoding response", err)
	}
}

func (t Transport) ShowStaticSinglePageApplication(w http.ResponseWriter, r *http.Request) {
	carsResp := t.getCarsList()
	tmpl, _ := template.ParseFiles("templates/index.html")
	err := tmpl.Execute(w, carsResp)
	if err != nil {
		println(err)
	}
}

func (t Transport) CreateCar(response http.ResponseWriter, request *http.Request) {
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

	if err := t.db.Create(&car).Error; err != nil {
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

func (t Transport) ShowCarByID(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "id")

	var carModel models.Car
	if err := t.db.First(&carModel, key).Error; err != nil {
		panic(err) //ToDo: handle error!!!
	}

	carDef := Car{
		ID:           carModel.ID,
		ModelName:    carModel.ModelName,
		Type:         carModel.Type,
		Transmission: carModel.Transmission,
		Engine:       carModel.Engine,
		HorsePower:   carModel.HorsePower,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(carDef)
	if err != nil {
		return
	}
}

func (t Transport) EditCarById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	requestBody, _ := ioutil.ReadAll(r.Body)
	var ca UpdateCarRequest

	logrus.WithField("editCarById", ca).Info("starting editCarById")

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

	if err := t.db.Where("id = ?", id).Updates(model); err == nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ca)
	if err != nil {
		return
	}
}

func (t Transport) DeleteCar(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "id")

	logrus.WithField("car_id", key).Info("starting deletion car")

	var car models.Car
	if err := t.db.Debug().Where("id = ?", key).Delete(&car).Error; err != nil {
		logrus.WithField("error", err).Error("error deleting car")
	}

	w.WriteHeader(http.StatusOK)
}

type InternalServerError struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
	Code    int    `json:"code"`
}

func FireInternalServerError(w http.ResponseWriter, message string, error error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	ierr := InternalServerError{Message: message, Error: error, Code: http.StatusInternalServerError}

	json.NewEncoder(w).Encode(ierr)
}
