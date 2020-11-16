package mysqldb

import (
	"github.com/artofimagination/mysql-user-db-go-interface/models"
)

// GetAllFeatures returns all fetures stored in the database.
func GetAllFeatures() ([]models.Feature, error) {
	queryString := "select id, name, config from features"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query, err := db.Query(queryString)

	if err != nil {
		return nil, err
	}

	defer query.Close()

	features := []models.Feature{}
	for query.Next() {
		var feature models.Feature
		err := query.Scan(&feature.ID, &feature.Name, &feature.Config)
		if err != nil {
			return nil, err
		}
		features = append(features, feature)
	}

	return features, nil
}
