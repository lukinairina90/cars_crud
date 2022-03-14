package car

import (
	"C_CRUD/configuration"
	"C_CRUD/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"html/template"
	"io/ioutil"
	"net/http"
)

func NewTransport(db *gorm.DB, cfg configuration.Config) *Transport {
	return &Transport{
		db:  db,
		cfg: cfg,
	}
}

type Transport struct {
	db  *gorm.DB
	cfg configuration.Config
}

func (t Transport) ShowCars(w http.ResponseWriter, r *http.Request) {
	carsResp := t.getCarsList()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(carsResp)
	if err != nil {
		return
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
	var ca CarUpdate

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

	w.WriteHeader(http.StatusNoContent)
}

func (t Transport) getCarsList() []Car {
	var carsModels []models.Car
	t.db.Find(&carsModels)

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

type CreateCarRequest struct {
	ModelType    string `json:"model_type"`
	Type         string `json:"type"`
	Transmission string `json:"transmission"`
	Engine       string `json:"engine"`
	HorsePower   string `json:"horse_power"`
}

type Car struct {
	ID        uint   `json:"id"`
	ModelName string `json:"model_name"`
	Type      string `json:"model_type"`

	//ModelInfo    string `json:"model_info"` // ModelName (Type)

	Transmission string `json:"transmission"`
	Engine       string `json:"engine"`
	HorsePower   string `json:"horse_power"`
}

type CarUpdate struct {
	ModelName    string `json:"model_name"`
	Type         string `json:"type"`
	Transmission string `json:"transmission"`
	Engine       string `json:"engine"`
	HorsePower   string `json:"horse_power"`
}
