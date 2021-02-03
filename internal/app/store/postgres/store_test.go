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
				RETURNING id
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

func TestGetOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := `
				SELECT a.name, a.price, a.description, p.url
				FROM ads AS a
				INNER JOIN photos AS p
				ON p.ad_id = a.id
				WHERE a.id = $1
			`

	expRows := mock.NewRows(
		[]string{"name", "price", "description", "url"},
	).AddRow("1", "4", "qwerty", "http://123/image1").
		AddRow("1", "5", "qwerty", "http://123/image2").
		AddRow("1", "3", "qwerty", "http://123/image3")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(expRows)
	mock.ExpectCommit()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := New(sqlxDB)

	_, err = store.Add().GetOne(1, true, true)
	assert.NoError(t, err)
}

func TestGetList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mainQuery := `
				SELECT s.name, s.price, s.url
				FROM
				(SELECT DISTINCT ON (a.id) a.id, a.name, a.price, p.url, a.created_at
				FROM ads as a
				INNER JOIN photos AS p
				ON p.ad_id = a.id) as s
				ORDER BY created_at asc
				LIMIT 10 OFFSET 0
	`

	subQuery := `
				SELECT count(*) FROM ads
	`
	expRows := mock.NewRows(
		[]string{"name", "price", "url"},
	).AddRow("ad1", "4", "http://123/image1").
		AddRow("ad2", "5", "http://123/image2").
		AddRow("ad3", "3", "http://123/image3")

	expRowsSubQuery := mock.NewRows(
		[]string{"count"},
	).AddRow("3")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(subQuery)).WillReturnRows(expRowsSubQuery)
	mock.ExpectQuery(regexp.QuoteMeta(mainQuery)).WithArgs(1).WillReturnRows(expRows)
	mock.ExpectCommit()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := New(sqlxDB)

	_, _, err = store.Add().GetList(10, 1, "created_at", "asc")
	assert.NoError(t, err)
}
