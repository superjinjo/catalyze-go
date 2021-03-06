package users

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//Model represents basic user data stored in database excluding the password
type Model struct {
	Username  string
	Firstname string
	Lastname  string
	Color     string
	ID        int
}

//Repository handles getting and storing user data from database
type Repository struct {
	DBuser string
	DBpw   string
	DBhost string
	DBname string
	conn   *sql.DB
}

//GetUser gets full user data based on username
func (repo *Repository) GetUser(username string) (Model, error) {
	var data Model

	db, err := repo.getConnection()
	if err != nil {
		return data, err
	}

	var (
		id                         int
		firstname, lastname, color string
	)

	row := db.QueryRow("SELECT id, firstname, lastname, color FROM users WHERE username = ?", username)
	if err := row.Scan(&id, &firstname, &lastname, &color); err != nil {
		return data, err
	}

	data = Model{username, firstname, lastname, color, id}

	return data, nil
}

//GetUserByID gets full user data based on user id
func (repo *Repository) GetUserByID(id int) (Model, error) {
	var data Model

	db, err := repo.getConnection()
	if err != nil {
		return data, err
	}

	var username, firstname, lastname, color string

	row := db.QueryRow("SELECT username, firstname, lastname, color FROM users WHERE id = ?", id)
	if err := row.Scan(&username, &firstname, &lastname, &color); err != nil {
		return data, err
	}

	data = Model{username, firstname, lastname, color, id}

	return data, nil
}

//GetUsername gets the username associated with the given user id
func (repo *Repository) GetUsername(id int) (string, error) {
	var username string

	db, err := repo.getConnection()
	if err != nil {
		return "", err
	}

	row := db.QueryRow("SELECT username FROM users WHERE id = ?", id)
	if err := row.Scan(&username); err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("User not found")
		}
		return "", err
	}

	return username, nil
}

//CheckCredentials checks the given username and password to see if they match in the database
func (repo *Repository) CheckCredentials(username, password string) (int, error) {
	var (
		id       int
		hashedpw string
	)

	db, err := repo.getConnection()
	if err != nil {
		return 0, err
	}

	row := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username)
	if err := row.Scan(&id, &hashedpw); err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("Username and password did not match")
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedpw), []byte(password)); err == nil {
		return id, nil
	} else {
		return 0, errors.New("Username and password did not match")
	}

}

//InsertUser adds a user to the database and hashes the password
func (repo *Repository) InsertUser(username, password, firstname, lastname, color string) (Model, error) {
	var data Model

	db, err := repo.getConnection()

	if err != nil {
		return data, err
	}

	hashedpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return data, err
	}

	stmt, err := db.Prepare("INSERT INTO users(username, password, firstname, lastname, color) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return data, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(username, string(hashedpw), firstname, lastname, color)
	if err != nil {
		return data, err
	}

	if id, err := result.LastInsertId(); err == nil {
		data = Model{username, firstname, lastname, color, int(id)}
		return data, nil
	}

	return data, err
}

//UpdateUser updates the name and favorite color of the user with the given username
func (repo *Repository) UpdateUser(username, firstname, lastname, color string) (Model, error) {
	var data Model

	db, err := repo.getConnection()
	if err != nil {
		return data, err
	}

	stmt, err := db.Prepare("UPDATE users SET firstname = ?, lastname = ?, color = ? WHERE username = ?")
	if err != nil {
		return data, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(firstname, lastname, color, username)
	if err != nil {
		return data, err
	}

	return repo.GetUser(username)
}

//DeleteUser removes user from the database
func (repo *Repository) DeleteUser(username string) error {
	db, err := repo.getConnection()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE username = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(username)

	return err
}

func (repo *Repository) buildDsn() string {
	dsnslice := []string{repo.DBuser, ":", repo.DBpw, "@tcp(", repo.DBhost, ")/", repo.DBname}

	return strings.Join(dsnslice, "")
}

func (repo *Repository) getConnection() (*sql.DB, error) {
	var err error

	if repo.conn == nil {
		repo.conn, err = sql.Open("mysql", repo.buildDsn())
	}

	return repo.conn, err
}

//Close closes database connection
func (repo *Repository) Close() error {
	var err error
	if repo.conn != nil {
		err = repo.conn.Close()
	}

	return err
}
