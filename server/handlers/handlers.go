package handlers

import (
	"database/sql"
	"encoding/json"
	"history-map/server/db"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func GetAllMaps(w http.ResponseWriter, r *http.Request) {
	rows, err := db.GetDB().Query("SELECT id, name, year, image_path, ST_AsGeoJSON(bounds) FROM historical_maps")
	if err != nil {
		log.Println(err)
		http.Error(w, "Database query error.", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	maps := []db.HistoricalMap{}
	for rows.Next() {
		var m db.HistoricalMap
		var bound string
		if err := rows.Scan(&m.ID, &m.Name, &m.Year, &m.ImagePath, &bound); err != nil {
			log.Println(err)
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}
		geoGem, err := geojson.UnmarshalGeometry([]byte(bound))
		if err != nil {
			log.Println("unmarshall error:", err)
			http.Error(w, "Error unmarshall bound", http.StatusInternalServerError)
			return
		}
		geom, ok := geoGem.Geometry().(orb.Polygon)
		if !ok {
			log.Println("geometry is not a polygon")
			http.Error(w, "Geometry is not a Polygon", http.StatusInternalServerError)
			return
		}
		m.Bounds = geom
		maps = append(maps, m)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		http.Error(w, "Row error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(maps); err != nil {
		log.Println(err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetMapByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid maps ID", http.StatusBadRequest)
		return
	}

	row := db.GetDB().QueryRow("SELECT id, name, year, image_path, ST_AsGeoJSON(bounds) FROM historical_maps WHERE id = $1", id)
	var m db.HistoricalMap
	var bound string
	if err := row.Scan(&m.ID, &m.Name, &m.Year, &m.ImagePath, &bound); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Map not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	geoGem, err := geojson.UnmarshalGeometry([]byte(bound))
	if err != nil {
		log.Println("unmarshall error", err)
		http.Error(w, "Error unmarshalling bound", http.StatusInternalServerError)
		return
	}
	geom, ok := geoGem.Geometry().(orb.Polygon)
	if !ok {
		log.Println("geometry is not a polygon")
		http.Error(w, "Geometry is not a Polygon", http.StatusInternalServerError)
		return
	}

	m.Bounds = geom

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Println(err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

func CreateMap(w http.ResponseWriter, r *http.Request) {
	var newMap db.HistoricalMap
	if err := json.NewDecoder(r.Body).Decode(&newMap); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gj := geojson.NewGeometry(newMap.Bounds)
	boundJSON, err := gj.MarshalJSON()
	if err != nil {
		http.Error(w, "Error marshalling bounds to JSON", http.StatusInternalServerError)
		return
	}

	query := `
        INSERT INTO historical_maps (name, year, image_path, bounds)
        VALUES ($1, $2, $3, ST_GeomFromGeoJSON($4))
        RETURNING id
    `

	var insertedID int
	err = db.GetDB().QueryRow(query, newMap.Name, newMap.Year, newMap.ImagePath, boundJSON).Scan(&insertedID)

	if err != nil {
		log.Println(err)
		http.Error(w, "Database insert error", http.StatusInternalServerError)
		return
	}
	newMap.ID = insertedID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newMap); err != nil {
		log.Println(err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
