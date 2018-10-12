package user

import (
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

type (
	UserModel interface {
		FindByID(id int) (*User, error)
		FindAll() ([]User, error)
		FindByEmail(email string) (*User, error)
		CreateUserByEmail(email string, passwordHash string) error
		CreateUserByFacebook(facebookId string) error
		SetPassword(userId int, newPassword string) error
	}

	UserModelImpl struct {
		db *sqlx.DB
	}

	User struct {
		ID           int         `json:"id" db:"id"`
		Email        null.String `json:"email" db:"email"`
		Password     null.String `json:"-"`
		CreationType string      `json:"creationType" db:"creation_type"`
		FacebookId   null.String `json:"facebookId" db:"facebook_id"`
	}
)

func NewUserModel(db *sqlx.DB) *UserModelImpl {
	return &UserModelImpl{
		db: db,
	}
}

func (u *UserModelImpl) FindByID(id int) (*User, error) {
	var users []User

	err := u.db.Select(&users, "SELECT * FROM users where id = $1 limit 1", id)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}

func (u *UserModelImpl) FindAll() ([]User, error) {
	var users []User
	e := u.db.Select(&users, "SELECT * FROM users order by id asc")
	if e != nil {
		log.Errorf("An error occurred during get users %v", e)
		return nil, e
	}
	return users, nil
}

func (u *UserModelImpl) FindByEmail(email string) (*User, error) {
	var users []User

	err := u.db.Select(&users, "SELECT * FROM users where email = $1 limit 1", email)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}

func (u *UserModelImpl) CreateUserByEmail(email, passwordHash string) error {
	_, err := u.db.Exec("INSERT INTO users (email, password, creation_type) VALUES ($1, $2, 'email')", email, passwordHash)
	return err
}

func (u *UserModelImpl) CreateUserByFacebook(facebookId string) error {
	_, err := u.db.Exec("INSERT INTO users (facebook_id, creation_type) VALUES ($1, 'facebook')", facebookId)
	return err
}

func (u *UserModelImpl) SetPassword(userId int, newPassword string) error {
	_, err := u.db.Exec("UPDATE users SET password = $1 WHERE id = %2", newPassword, userId)
	return err
}