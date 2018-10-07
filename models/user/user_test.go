package user

import (
	"testing"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

func MockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlx.DB) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expecting", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return mockDB, mock, sqlxDB
}


func TestFindByID(t *testing.T) {
	mockDB, mock, sqlxDB := MockDB(t)
	defer mockDB.Close()

	var cols []string = []string{"id", "name"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(1, "foobar"))

	um := NewUserModel(sqlxDB)
	u, e := um.FindByID("1")
	assert.Nil(t, e)

	expect := &User{
		ID:    1,
		Email: "foobar",
	}
	assert.Equal(t, expect, u)
}

func TestFindAll(t *testing.T) {
	mockDB, mock, sqlxDB := MockDB(t)
	defer mockDB.Close()

	u1 := User{ID: 1, Email: "foobar"}
	u2 := User{ID: 2, Email: "barbaz"}

	var cols []string = []string{"id", "name"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(u1.ID, u1.Email).
		AddRow(u2.ID, u2.Email))

	um := NewUserModel(sqlxDB)
	u, _ := um.FindAll()

	expect := []User{}
	expect = append(expect, u1, u2)
	assert.Equal(t, expect, u)
}
