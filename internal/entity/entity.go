package entity

type Vehicle struct {
	VIN    string `json:"vin" required:"true"`
	Brand  string `json:"brand"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Status string `json:"status"`
}

type User struct {
	Name        string `json:"username"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Description string `json:"password"`
}

func NewVehicle(vin string, brand string, model string, year int, status string) *Vehicle {
	/*
		status (статус) — состояние автомобиля. Возможные значения:

		available (доступен) — готов к назначению на рейс.
		on_route (в рейсе) — находится в пути /зарезервирован.
		under_maintenance (на ремонте) — проходит ТО или ремонт.
	*/
	return &Vehicle{
		VIN:    vin,
		Brand:  brand,
		Model:  model,
		Year:   year,
		Status: status,
	}
}
