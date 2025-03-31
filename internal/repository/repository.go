package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DblMOKRQ/introductory-practice/internal/entity"
)

type Repository struct {
	dbVehicle *sql.DB
	dbUser    *sql.DB
}

func NewRepository(dbVehicle *sql.DB, dbUser *sql.DB) *Repository {
	return &Repository{dbVehicle: dbVehicle, dbUser: dbUser}
}

func (r *Repository) AddVehicle(v *entity.Vehicle) error {

	_, err := r.dbVehicle.Exec("INSERT INTO vehicles (vin,brand,model,year,status) VALUES ($1,$2,$3,$4,$5)", v.VIN, v.Brand, v.Model, v.Year, v.Status)
	if err != nil {
		return fmt.Errorf("failed to add vehicle: %w", err)
	}
	return nil

}

func (r *Repository) GetVehicle(vin string) (*entity.Vehicle, error) {
	var v entity.Vehicle
	var busyUntil time.Time
	var id int
	err := r.dbVehicle.QueryRow("SELECT * FROM vehicles WHERE vin = $1", vin).Scan(&id, &v.VIN, &v.Brand, &v.Model, &v.Year, &v.Status, &busyUntil)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}
	return &v, nil
}

func (r *Repository) UpdateVehicle(v *entity.Vehicle) error {
	_, err := r.dbVehicle.Exec("UPDATE vehicles SET brand = $1, model = $2, year = $3, status = $4 WHERE vin = $5", v.Brand, v.Model, v.Year, v.Status, v.VIN)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	return nil
}

func (r *Repository) DeleteVehicle(vin string) error {
	_, err := r.dbVehicle.Exec("DELETE FROM vehicles WHERE vin = $1", vin)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}
	return nil
}

func (r *Repository) GetAllVehicles() ([]*entity.Vehicle, error) {
	rows, err := r.dbVehicle.Query("SELECT * FROM vehicles")
	if err != nil {
		return nil, fmt.Errorf("failed to get all vehicles: %w", err)
	}
	defer rows.Close()
	var id int
	var vehicles []*entity.Vehicle
	for rows.Next() {
		var v entity.Vehicle
		var timeStamp time.Time
		err := rows.Scan(&id, &v.VIN, &v.Brand, &v.Model, &v.Year, &v.Status, &timeStamp)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, &v)
	}
	return vehicles, nil
}

func (r *Repository) UpdateStatus(vin string, status string) error {
	_, err := r.dbVehicle.Exec("UPDATE vehicles SET status = $1 WHERE vin = $2", status, vin)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

func (r *Repository) RentVehicle(vin string, u *entity.User) error {

	_, err := r.dbUser.Exec("INSERT INTO users (name,email,phone,reserve_vin,description) VALUES ($1,$2,$3,$4,$5)", u.Name, u.Email, u.Phone, vin, u.Description)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	err = r.UpdateStatus(vin, "on_route")
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}
