package auth

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)

const dateformat = "2006-01-02 15:04:05"

//Repository handles getting and storing user data from database
type Repository struct {
	DBuser string
	DBpw   string
	DBhost string
	DBname string
	conn   *sql.DB
}

//CreateToken generates an auth token for the given user
func (repo *Repository) CreateToken(userID, lifespan int) (string, error) {
	db, err := repo.getConnection()

	if err != nil {
		return "", err
	}

	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	now := time.Now()
	created := now.Format(dateformat)
	expires := now.Add(time.Second * time.Duration(lifespan)).Format(dateformat)

	stmt, err := db.Prepare("INSERT INTO auth(token, user_id, created, expires) VALUES(?, ?, ?, ?)")
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	_, err = stmt.Exec(token.String(), userID, created, expires)
	if err != nil {
		return "", err
	}

	return token.String(), err
}

//CheckToken checks to see if the token is valid and is not expired and returns the associated user id
func (repo *Repository) CheckToken(token string) (int, error) {
	db, err := repo.getConnection()
	if err != nil {
		return 0, err
	}

	var (
		userID  int
		expires string
	)

	row := db.QueryRow("SELECT user_id, expires FROM auth WHERE token = ?", token)
	if err := row.Scan(&userID, &expires); err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("Invalid token")
		}
		return 0, err
	}

	exp, _ := time.Parse(dateformat, expires)
	if exp.Before(time.Now()) {
		repo.DeleteToken(token)
		return 0, errors.New("Token expired")
	}

	return userID, nil
}

//DeleteToken removes token from database, thus invalidating it
func (repo *Repository) DeleteToken(token string) error {
	db, err := repo.getConnection()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("DELETE FROM auth WHERE token = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(token)

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
