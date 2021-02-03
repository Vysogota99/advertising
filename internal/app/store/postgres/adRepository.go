package postgres

import (
	"fmt"

	"github.com/Vysogota99/advertising/internal/app/models"
)

// AdRepository - postgres implementation of AddRepository
type AdRepository struct {
	store *Store
}

// Create ...
func (a *AdRepository) Create(ad models.Ad) (int, error) {
	tx, err := a.store.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		return 0, nil
	}

	query := `
		INSERT INTO ads(name, description, price)
		VALUES ($1, $2, $3) 
		RETURNING id
	`

	var id int
	if err = tx.QueryRowx(query, ad.Name, ad.Description, ad.Price).Scan(&id); err != nil {
		return 0, err
	}

	query = `
		INSERT INTO photos(ad_id, url)
		VALUES %s
	`

	n := len(ad.Links)
	valuesToQuery := "%s"

	params := []interface{}{}
	params = append(params, id)
	for i := 0; i < n; i++ {

		var inputValue string
		if i != n-1 {
			inputValue = fmt.Sprintf("($1, $%d), %%s", i+2)
		} else {
			inputValue = fmt.Sprintf("($1, $%d)", i+2)
		}

		valuesToQuery = fmt.Sprintf(valuesToQuery, inputValue)
		params = append(params, ad.Links[i])
	}

	query = fmt.Sprintf(query, valuesToQuery)

	_, err = tx.Exec(query, params...)
	if err != nil {
		return 0, nil
	}

	tx.Commit()
	return id, nil
}

// GetOne ...
func (a *AdRepository) GetOne(id int, description, photos bool) (*models.Ad, error) {
	tx, err := a.store.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		return nil, nil
	}

	query := `
		SELECT a.name, a.price, %s p.url
		FROM ads AS a
		INNER JOIN photos AS p
		ON p.ad_id = a.id
		WHERE a.id = $1
	`

	descriptionImputer := ""
	if description {
		descriptionImputer = "a.description,"
	}

	query = fmt.Sprintf(query, descriptionImputer)
	rows, err := tx.Queryx(query, id)
	if err != nil {
		return nil, err
	}

	ad := models.Ad{}
	links := []string{}

	for rows.Next() {
		if description {
			if err := rows.Scan(&ad.Name, &ad.Price, &ad.Description, &ad.Link); err != nil {
				return nil, err
			}
		} else {
			if err := rows.Scan(&ad.Name, &ad.Price, &ad.Link); err != nil {
				return nil, err
			}
		}

		links = append(links, ad.Link)
	}

	if len(links) == 0 {
		return nil, nil
	}

	if photos {
		ad.Links = links
	} else {
		ad.Links = []string{links[0]}
	}

	tx.Commit()
	return &ad, nil
}

// GetList ...
func (a *AdRepository) GetList(limit int, offset int, sortBy, sortDirection string) ([]models.Ad, int, error) {
	tx, err := a.store.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		return nil, 0, nil
	}

	query := `
		SELECT count(*) FROM ads
	`

	var nRows int
	if err := tx.QueryRow(query).Scan(&nRows); err != nil {
		return nil, 0, err
	}

	var nPages int

	if nRows%limit > 0 {
		nPages = nRows/limit + 1
	} else {
		nPages = nRows / limit
	}

	offset = limit * (offset - 1)
	if sortBy == "" {
		sortBy = "created_at"
	}

	if sortDirection == "" {
		sortDirection = "asc"
	}

	query = fmt.Sprintf(`
			SELECT s.name, s.price, s.url
			FROM
			(SELECT DISTINCT ON (a.id) a.id, a.name, a.price, p.url, a.created_at
			FROM ads as a
			INNER JOIN photos AS p
			ON p.ad_id = a.id) as s
			ORDER BY %s %s
			LIMIT $1 OFFSET $2
	`, sortBy, sortDirection)

	rows, err := tx.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := []models.Ad{}
	for rows.Next() {
		ad := models.Ad{}
		var url string
		if err := rows.Scan(&ad.Name, &ad.Price, &url); err != nil {
			return nil, 0, err
		}

		ad.Links = append(ad.Links, url)
		result = append(result, ad)
	}

	return result, nPages, nil
}
