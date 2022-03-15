package car

type CreateCarRequest struct {
	ModelType    string `json:"model_type"`
	Type         string `json:"type"`
	Transmission string `json:"transmission"`
	Engine       string `json:"engine"`
	HorsePower   string `json:"horse_power"`
}

type Car struct {
	ID           uint   `json:"id"`
	ModelName    string `json:"model_name"`
	Type         string `json:"model_type"`
	Transmission string `json:"transmission"`
	Engine       string `json:"engine"`
	HorsePower   string `json:"horse_power"`
}

type UpdateCarRequest struct {
	ModelName    string `json:"model_name"`
	Type         string `json:"type"`
	Transmission string `json:"transmission"`
	Engine       string `json:"engine"`
	HorsePower   string `json:"horse_power"`
}
