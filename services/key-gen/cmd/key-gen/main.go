package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "username"
	password = "password"
	dbname   = "url_shortener"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func base62Encode(num int64) string {
	if num == 0 {
		return string(alphabet[0])
	}
	chars := []string{}
	base := int64(len(alphabet))
	for num > 0 {
		rem := num % base
		chars = append([]string{string(alphabet[rem])}, chars...)
		num = num / base
	}
	return strings.Join(chars, "")
}

func incrementCounter(db *sql.DB) (int64, error) {
	var newValue int64
	err := db.QueryRow(
		"UPDATE counter SET value = value + 1 WHERE id = 1 RETURNING value",
	).Scan(&newValue)

	if err != nil {
		return 0, err
	}
	return newValue, nil
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	counterValue, err := incrementCounter(db)
	if err != nil {
		log.Fatal(err)
	}

	urlAlias := base62Encode(counterValue)
	fmt.Println("Generated URL Alias:", urlAlias)
}
