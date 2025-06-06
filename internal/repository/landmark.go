package repository

import (
	"context"
	"database/sql"
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
			st_astext(landmark.location) as loc,
			landmark.images_name
		FROM landmark        
		JOIN public.weather w on landmark.id = w.landmark_id
		
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
		err := rows.Scan(&f.ID, &f.Name, &f.Address, &p, &f.ImagePath)
		if err != nil {
			return []models.Landmark{}, err
		}

		f.TranslatedName = strings.Split(f.ImagePath, ".")[0]
		f.Location = utils.LocationFromPoint(p)
		facilities = append(facilities, f)
	}
	return facilities, nil

}

func (l *LandmarkDB) GetLandmarks(page int, categories []string) ([]models.Landmark, error) {
	offset := (page - 1) * PageSize
	query := `
			SELECT
				landmark.id,
				landmark.name,
				landmark.address,
				landmark.category,
				landmark.description,
				landmark.history,
				st_astext(landmark.location) as loc		,
				landmark.images_name
			    FROM landmark
			
		`

	if len(categories) > 0 {
		query = fmt.Sprintf(`
		SELECT
			landmark.id,
			landmark.name,
			landmark.address,
			landmark.category,
			landmark.description,
			landmark.history,
			st_astext(landmark.location) AS loc,
			landmark.images_name
		FROM landmark
		
		WHERE lower(landmark.category) IN (%s)
		
	`, strings.Join(categories, ","))
	}
	var rows *sql.Rows
	var err error
	if page != -1 {
		query += " LIMIT $1 OFFSET $2"
		rows, err = l.postgres.Query(query, PageSize, offset)
	} else {
		rows, err = l.postgres.Query(query)
	}
	if err != nil {
		return []models.Landmark{}, err
	}
	defer rows.Close()
	var landmarks []models.Landmark
	for rows.Next() {
		f, p := models.Landmark{}, ""
		err := rows.Scan(&f.ID, &f.Name, &f.Address, &f.Category, &f.Description, &f.History, &p, &f.ImagePath)
		if err != nil {
			return []models.Landmark{}, err

		}
		f.TranslatedName = strings.Split(f.ImagePath, ".")[0]

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
			ST_AsText(landmark.location) AS loc,
			landmark.images_name
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
		err := res.Scan(&f.ID, &f.Name, &f.Address, &f.Category, &f.Description, &f.History, &p, &f.ImagePath)
		if err != nil {
			return []models.Landmark{}, err

		}

		f.TranslatedName = strings.Split(f.ImagePath, ".")[0]
		f.Location = utils.LocationFromPoint(p)
		landmarks = append(landmarks, f)
	}
	return landmarks, nil
}

func (l *LandmarkDB) Search(q string) ([]models.Landmark, error) {
	landmarksMap := make(map[int]models.Landmark)

	query := `
		SELECT
			landmark.id,
			landmark.name,
			landmark.address,
			landmark.category,
			landmark.description,
			landmark.history,
			st_astext(landmark.location) as loc,
			landmark.images_name
		FROM landmark
		WHERE to_tsvector('russian', landmark.name) @@ to_tsquery('russian', $1)
	`
	rows, err := l.postgres.Query(query, q)
	if err != nil {
		return []models.Landmark{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Landmark
		var p string
		err := rows.Scan(&f.ID, &f.Name, &f.Address, &f.Category, &f.Description, &f.History, &p, &f.ImagePath)
		if err != nil {
			return []models.Landmark{}, err
		}
		f.TranslatedName = strings.Split(f.ImagePath, ".")[0]
		f.Location = utils.LocationFromPoint(p)
		landmarksMap[f.ID] = f
	}
	if err = rows.Err(); err != nil {
		return []models.Landmark{}, err
	}

	query = `
		SELECT
			landmark.id,
			landmark.name,
			landmark.address,
			landmark.category,
			landmark.description,
			landmark.history,
			st_astext(landmark.location) as loc,
			landmark.images_name
		FROM landmark
		WHERE to_tsvector('russian', landmark.address) @@ to_tsquery('russian', $1)
	`
	rows, err = l.postgres.Query(query, q)
	if err != nil {
		return []models.Landmark{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Landmark
		var p string
		err := rows.Scan(&f.ID, &f.Name, &f.Address, &f.Category, &f.Description, &f.History, &p, &f.ImagePath)
		if err != nil {
			return []models.Landmark{}, err
		}
		_, ok := landmarksMap[f.ID]
		if ok {
			continue
		}

		f.TranslatedName = strings.Split(f.ImagePath, ".")[0]
		f.Location = utils.LocationFromPoint(p)
		landmarksMap[f.ID] = f
	}
	if err = rows.Err(); err != nil {
		return []models.Landmark{}, err
	}

	landmarks := make([]models.Landmark, 0, len(landmarksMap))
	for _, landmark := range landmarksMap {
		landmarks = append(landmarks, landmark)
	}

	return landmarks, nil
}

func (l *LandmarkDB) UpdateImagePath(place string, path string) error {
	query :=
		`
		UPDATE landmark SET images_name = $1 WHERE name = $2 or lower(images_name)=$1
		`
	_, err := l.postgres.Exec(query, path, place)
	return err

}
func (l *LandmarkDB) GetLandmarksByName(name string) (models.Landmark, error) {
	query :=
		`
		SELECT id,name, address,category,description,history,st_astext(landmark.location) as loc,images_name FROM landmark
		WHERE SUBSTRING(images_name FROM 1 FOR POSITION('.' IN images_name) ) LIKE $1 || '%';
		`
	res := l.postgres.QueryRow(query, name)
	landmark := models.Landmark{}
	var tmp string
	if err := res.Scan(&landmark.ID, &landmark.Name, &landmark.Address, &landmark.Category, &landmark.Description, &landmark.History, &tmp, &landmark.ImagePath); err != nil {
		return models.Landmark{}, err
	}
	landmark.TranslatedName = strings.Split(landmark.ImagePath, ".")[0]
	landmark.Location = utils.LocationFromPoint(tmp)
	return landmark, nil
}

func (l *LandmarkDB) GetLandmarksByCategories(categories []string) ([]models.Landmark, error) {
	if len(categories) == 0 {
		return []models.Landmark{}, nil
	}

	placeholders := make([]string, len(categories))
	args := make([]interface{}, len(categories))
	for i, category := range categories {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = category
	}

	query := fmt.Sprintf(`
		SELECT
			landmark.id,
			landmark.name,
			landmark.address,
			landmark.category,
			landmark.description,
			landmark.history,
			st_astext(landmark.location) AS loc,
			landmark.images_name
		FROM landmark
		WHERE landmark.category IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := l.postgres.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var landmarks []models.Landmark
	for rows.Next() {
		var f models.Landmark
		var loc string
		if err := rows.Scan(&f.ID, &f.Name, &f.Address, &f.Category, &f.Description, &f.History, &loc, &f.ImagePath); err != nil {
			return nil, err
		}
		f.TranslatedName = strings.Split(f.ImagePath, ".")[0]
		f.Location = utils.LocationFromPoint(loc)
		landmarks = append(landmarks, f)
	}

	return landmarks, nil
}
