package database

import (
	"database/sql"
	"log/slog"
)

func (s *Database) GetStyluses() ([]Stylus, error) {
	query := `SELECT * FROM styluses ORDER BY name`
	rows, err := s.DB.Query(query)
	if err != nil {
		slog.Error("Failed to get styluses", "error", err)
		return nil, err
	}
	defer rows.Close()

	var styluses []Stylus
	for rows.Next() {
		var stylus Stylus
		var purchaseDate sql.NullTime

		err := rows.Scan(
			&stylus.ID,
			&stylus.Name,
			&stylus.Manufacturer,
			&stylus.ExpectedLifespan,
			&purchaseDate,
			&stylus.Active,
			&stylus.Primary,
			&stylus.ModelNumber,
			&stylus.CreatedAt,
			&stylus.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan stylus", "error", err)
			return nil, err
		}

		if purchaseDate.Valid {
			stylus.PurchaseDate = &purchaseDate.Time
		}

		styluses = append(styluses, stylus)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating stylus rows", "error", err)
		return styluses, err
	}

	return styluses, nil
}

// INFO: Not used, left for reference
// GetStylusByID retrieves a stylus by its ID
func (s *Database) GetStylusByID(id int) (*Stylus, error) {
	query := `SELECT * FROM styluses WHERE id = ?`
	row := s.DB.QueryRow(query, id)

	var stylus Stylus
	var purchaseDate sql.NullTime

	err := row.Scan(
		&stylus.ID,
		&stylus.Name,
		&stylus.Manufacturer,
		&stylus.ExpectedLifespan,
		&purchaseDate,
		&stylus.Active,
		&stylus.Primary,
		&stylus.ModelNumber,
		&stylus.CreatedAt,
		&stylus.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No stylus found
		}
		slog.Error("Failed to scan stylus", "error", err)
		return nil, err
	}

	if purchaseDate.Valid {
		stylus.PurchaseDate = &purchaseDate.Time
	}

	return &stylus, nil
}

func (s *Database) CreateStylus(stylus *Stylus) error {
	query := `
		INSERT INTO styluses (
			name, manufacturer, expected_lifespan_hours, purchase_date, 
			active, primary_stylus, model_number
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	var purchaseDate interface{}
	if stylus.PurchaseDate != nil {
		purchaseDate = stylus.PurchaseDate
	} else {
		purchaseDate = nil
	}

	err := s.DB.QueryRow(
		query,
		stylus.Name,
		stylus.Manufacturer,
		stylus.ExpectedLifespan,
		purchaseDate,
		stylus.Active,
		stylus.Primary,
		stylus.ModelNumber,
	).Scan(&stylus.ID, &stylus.CreatedAt, &stylus.UpdatedAt)
	if err != nil {
		slog.Error("Failed to create stylus", "error", err)
		return err
	}

	return nil
}

func (s *Database) UpdateStylus(stylus *Stylus) error {
	query := `
		UPDATE styluses SET
			name = ?,
			manufacturer = ?,
			expected_lifespan_hours = ?,
			purchase_date = ?,
			active = ?,
			primary_stylus = ?,
			model_number = ?
		WHERE id = ?
		RETURNING updated_at
	`

	var purchaseDate interface{}
	if stylus.PurchaseDate != nil {
		purchaseDate = stylus.PurchaseDate
	} else {
		purchaseDate = nil
	}

	err := s.DB.QueryRow(
		query,
		stylus.Name,
		stylus.Manufacturer,
		stylus.ExpectedLifespan,
		purchaseDate,
		stylus.Active,
		stylus.Primary,
		stylus.ModelNumber,
		stylus.ID,
	).Scan(&stylus.UpdatedAt)
	if err != nil {
		slog.Error("Failed to update stylus", "error", err)
		return err
	}

	return nil
}

func (s *Database) DeleteStylus(id int) error {
	// First check if this stylus is referenced in play_history
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM play_history WHERE stylus_id = ?", id).Scan(&count)
	if err != nil {
		slog.Error("Failed to check play history references", "error", err)
		return err
	}

	if count > 0 {
		// If referenced, update play_history to set stylus_id to NULL
		_, err = s.DB.Exec("UPDATE play_history SET stylus_id = NULL WHERE stylus_id = ?", id)
		if err != nil {
			slog.Error("Failed to update play history references", "error", err)
			return err
		}
	}

	// Now delete the stylus
	_, err = s.DB.Exec("DELETE FROM styluses WHERE id = ?", id)
	if err != nil {
		slog.Error("Failed to delete stylus", "error", err)
		return err
	}

	return nil
}
