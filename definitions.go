package main

type Car struct {
	ID        uint   `json:"id"`
	ModelName string `json:"model_name"`
	Type      string `json:"model_type"`

	//ModelInfo    string `json:"model_info"` // ModelName (Type)

	Transmission string `json:"transmission"`
	Engine       string `json:"engine"`
	HorsePower   string `json:"horse_power"`
}
