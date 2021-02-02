package postgres

import (
	"fmt"
	"log"

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
		VALUES ($1, $2, $3) RETURNING id
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
		log.Println(err)
		return 0, nil
	}

	tx.Commit()
	return 1, nil
}

// GetOne ...
func (a *AdRepository) GetOne(int) (models.Ad, error) {
	result := models.Ad{}
	return result, nil
}

// GetList ...
func (a *AdRepository) GetList(curr int, limit int, offset int) ([]models.Ad, error) {
	return nil, nil
}
