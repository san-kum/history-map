package handlers

import (
	"encoding/json"
	"history-map/server/db"
	"log"
	"net/http"
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
		if err := rows.Scan(&m.ID, &m.Name, &m.ImagePath, &bound); err != nil {
			log.Println(err)
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal([]byte(bound), &m.Bounds); err != nil {
			log.Println("unmarshall error:", err)
			http.Error(w, "Error unmarshall bound", http.StatusInternalServerError)
		}
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
