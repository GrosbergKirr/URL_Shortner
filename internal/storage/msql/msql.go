package mysql

import (
	"awesomeProject/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.mysql.New"

	db, err := sql.Open("mysql", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//defer db.Close()

	a, err := db.Exec("CREATE TABLE IF NOT EXISTS url (id INTEGER PRIMARY KEY, alias varchar(8) NOT NULL UNIQUE, url TEXT NOT NULL, UNIQUE (alias))")
	if err != nil {
		panic(err)
	}
	_ = a
	if err != nil {
		return nil, fmt.Errorf("{op}: #{err}")
	}

	return &Storage{db: db}, nil
}
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.msql.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url,alias) values(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s : %w", op, err)
	}
	_ = res
	return 0, nil
	//id, err := res.LastInsertId()
	//if err != nil {
	//	return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	//}
	//
	//return id, nil

}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.msql.GetUrl"

	stmt, err := s.db.Prepare("select url from url where alias = (?)")
	if err != nil {
		return "0", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}
