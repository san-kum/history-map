package db

import "github.com/paulmach/orb"

type HistoricalMap struct {
	ID        int         `json:"id" db:"id"`
	Name      string      `json:"name" db:"name"`
	Year      int         `json:"year" db:"year"`
	ImagePath string      `json:"image_path" db:"image_path"`
	Bounds    orb.Polygon `json:"bounds" db:"bounds"`
}
