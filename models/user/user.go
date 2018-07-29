package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"strconv"
)

type (
	UserModel interface {
		FindByID(id string) (*User, error)
		FindAll() ([]User, error)
		FindByLogin(login string) (*User, error)
		CreateUser(login string, passwordHash string) (error)
	}

	UserModelImpl struct {
		db *sqlx.DB
	}

	User struct {
		ID   int    `json:"id" db:"id"`
		Name string `json:"name" db:"name"`
		Password string `json:"-"`
		Surname string
		Lastname string
	}

)

func NewUserModel(db *sqlx.DB) *UserModelImpl {
	return &UserModelImpl{
		db: db,
	}
}

func (u *UserModelImpl) FindByID(id string) (*User, error) {
	var users []User

	err := u.db.Select(&users, "SELECT * FROM users where id = $1 limit 1", id)
	if err != nil {
		return nil, err
	}
	if len(users)==0 {
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

func (u *UserModelImpl) FindByLogin(login string) (*User, error) {
	var users []User

	err := u.db.Select(&users, "SELECT * FROM users where name = $1 limit 1", login)
	if err != nil {
		return nil, err
	}
	if len(users)==0 {
		return nil, nil
	}
	return &users[0], nil
}

func (u *UserModelImpl) CreateUser(login, passwordHash string) error {
	_, error := u.db.Exec("INSERT INTO users (name, surname, lastname, password) VALUES ($1, '', '', $2)", login, passwordHash)
	return error
}

func (u *User) GetPID() (pid string){
	return strconv.Itoa(u.ID)
}

func (u *User) PutPID(pid string) {
	if id, e := strconv.Atoi(pid); e != nil {
		log.Panicf("Cannot convert to")
	} else {
		u.ID = id
	}
}

func (u *User) GetPassword() (password string){
	return u.Password
}

func (u *User) PutPassword(password string){
	u.Password = password
}