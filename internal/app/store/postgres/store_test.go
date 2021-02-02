package postgres

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vysogota99/advertising/internal/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := `
				INSERT INTO ads(name, description, price)
				VALUES ($1, $2, $3)
			`

	ad := models.Ad{
		Name:        "iphone1",
		Description: "Cool description",
		Price:       1999.99,
		Links: []string{
			"http://images/1",
			"http://images/2",
			"http://images/3",
		},
	}

	expID := mock.NewRows(
		[]string{"id"},
	).AddRow("1")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(ad.Name, ad.Description, ad.Price).WillReturnRows(expID)

	query = `
			INSERT INTO photos(ad_id, url)
			VALUES ($1, $2), ($1, $3), ($1, $4)
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(1, "http://images/1", "http://images/2", "http://images/3").WillReturnResult(sqlmock.NewResult(3, 3))
	mock.ExpectCommit()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := New(sqlxDB)

	n, err := store.Add().Create(ad)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}
