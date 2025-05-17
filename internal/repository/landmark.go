package repository

import (
	"context"
	"fmt"
	"strings"

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

func (l *LandmarkDB) GetLandmarksByIDs(ids []any) ([]models.Landmark, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT
			landmark.id,
			landmark.name,
			landmark.address,
			landmark.category,
			landmark.description,
			landmark.history,
			ST_AsText(landmark.location) AS loc
		FROM landmark
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))
	res, err := l.postgres.Query(query, ids...)
	if err != nil {
		return []models.Landmark{}, err
	}
	defer res.Close()

	var landmarks []models.Landmark
	for res.Next() {
		f, p := models.Landmark{}, ""
		err := res.Scan(&f.ID, &f.Name, &f.Address, &f.Category, &f.Description, &f.History, &p)
		if err != nil {
			return []models.Landmark{}, err

		}
		f.Location = utils.LocationFromPoint(p)
		landmarks = append(landmarks, f)
	}
	return landmarks, nil
}
