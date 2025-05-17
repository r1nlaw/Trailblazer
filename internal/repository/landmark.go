package repository

import (
	"context"

	"trailblazer/internal/models"
	"trailblazer/internal/utils"

	"github.com/jmoiron/sqlx"
)

type LandmarkDB struct {
	ctx      context.Context
	postgres *sqlx.DB
}

var PageSize = 10

func NewLandmarkPostgres(ctx context.Context, db *sqlx.DB) *LandmarkDB {
	return &LandmarkDB{ctx: ctx, postgres: db}
}

func (l *LandmarkDB) GetFacilities(bbox models.BBOX) ([]models.Landmark, error) {
	const selectQuery = `
		SELECT 
			landmark.id,
			landmark.name,
			landmark.address,
			st_astext(landmark.location) as loc
		FROM landmark        
	  	WHERE ST_Intersects(ST_MakeEnvelope($1,$2,$3,$4,4326 ), landmark.location::geometry)
	  `
	rows, err := l.postgres.Query(selectQuery, bbox.SW.Lng, bbox.SW.Lat, bbox.NE.Lng, bbox.NE.Lat)
	if err != nil {
		return []models.Landmark{}, err
	}

	defer rows.Close()
	var facilities []models.Landmark
	for rows.Next() {
		f, p := models.Landmark{}, ""
		err := rows.Scan(&f.ID, &f.Name, &f.Address, &p)
		if err != nil {
			return []models.Landmark{}, err
		}
		f.Location = utils.LocationFromPoint(p)
		facilities = append(facilities, f)
	}
	return facilities, nil

}

func (l *LandmarkDB) GetLandmarks(page int) ([]models.Landmark, error) {
	offset := (page - 1) * PageSize
	query := `
			SELECT
				landmark.id,
				landmark.name,
				landmark.address,
				landmark.category,
				landmark.description,
				landmark.history,
				st_astext(landmark.location) as loc				
			    FROM landmark
			LIMIT $1 OFFSET $2
		`
	rows, err := l.postgres.Query(query, PageSize, offset)
	if err != nil {
		return []models.Landmark{}, err

	}
	defer rows.Close()
	var landmarks []models.Landmark
	for rows.Next() {
		f, p := models.Landmark{}, ""
		err := rows.Scan(&f.ID, &f.Name, &f.Address, &f.Category, &f.Description, &f.History, &p)
		if err != nil {
			return []models.Landmark{}, err

		}
		f.Location = utils.LocationFromPoint(p)
		landmarks = append(landmarks, f)
	}
	return landmarks, nil
}
