package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"monk-commerce/models"
)

type mysqlCouponRepo struct {
	db *sql.DB
}

func NewMySQLCouponRepo(db *sql.DB) CouponRepo {
	query := `
	CREATE TABLE IF NOT EXISTS coupons (
		id INT AUTO_INCREMENT PRIMARY KEY,
		type VARCHAR(50) NOT NULL,
		details JSON NOT NULL,
		expiration_date DATETIME NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	return &mysqlCouponRepo{db: db}
}

func (r *mysqlCouponRepo) Create(coupon *models.Coupon) error {
	detailsJSON, err := json.Marshal(coupon.Details)
	if err != nil {
		return err
	}

	query := `INSERT INTO coupons (type, details, expiration_date) VALUES (?, ?, ?)`
	res, err := r.db.Exec(query, coupon.Type, string(detailsJSON), coupon.ExpirationDate)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	coupon.ID = int(id)
	return nil
}

func (r *mysqlCouponRepo) GetAll() ([]models.Coupon, error) {
	query := `SELECT id, type, details, expiration_date FROM coupons`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []models.Coupon
	for rows.Next() {
		var c models.Coupon
		var detailsStr string
		var expDate sql.NullTime

		err := rows.Scan(&c.ID, &c.Type, &detailsStr, &expDate)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(detailsStr), &c.Details); err != nil {
			return nil, err
		}

		if expDate.Valid {
			c.ExpirationDate = &expDate.Time
		}

		coupons = append(coupons, c)
	}
	return coupons, nil
}

func (r *mysqlCouponRepo) GetByID(id int) (models.Coupon, error) {
	query := `SELECT id, type, details, expiration_date FROM coupons WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var c models.Coupon
	var detailsStr string
	var expDate sql.NullTime

	err := row.Scan(&c.ID, &c.Type, &detailsStr, &expDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, ErrCouponNotFound
		}
		return c, err
	}

	if err := json.Unmarshal([]byte(detailsStr), &c.Details); err != nil {
		return c, err
	}

	if expDate.Valid {
		c.ExpirationDate = &expDate.Time
	}

	if c.ExpirationDate != nil && time.Now().After(*c.ExpirationDate) {
		return models.Coupon{}, ErrCouponNotFound
	}

	return c, nil
}

func (r *mysqlCouponRepo) Update(id int, coupon *models.Coupon) error {
	detailsJSON, err := json.Marshal(coupon.Details)
	if err != nil {
		return err
	}

	query := `UPDATE coupons SET type = ?, details = ?, expiration_date = ? WHERE id = ?`
	res, err := r.db.Exec(query, coupon.Type, string(detailsJSON), coupon.ExpirationDate, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCouponNotFound
	}
	coupon.ID = id
	return nil
}

func (r *mysqlCouponRepo) Delete(id int) error {
	query := `DELETE FROM coupons WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCouponNotFound
	}
	return nil
}
