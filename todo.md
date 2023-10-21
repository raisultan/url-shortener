[+] Deletion Endpoint
[+] Deletion Endpoint Tests
[ ] KGS
    [+] Caching Mechanism
        [+] Saves Newly Generated Alias For 24 Hours into Cache
        [+] On Read First Checks Cache
        [ ] Unit Tests
[ ] Analytics Storage

## KGS

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

// base62 characters
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// base62Encode encodes a number to base62
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

// incrementCounter increments the counter in the database and returns the new value
func incrementCounter(db *sql.DB) (int64, error) {
var newValue int64
err := db.QueryRow("UPDATE counter SET value = value + 1 WHERE id = 1 RETURNING value").Scan(&newValue)
if err != nil {
return 0, err
}
return newValue, nil
}

func main() {
// Setup database connection
psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
host, port, user, password, dbname)
db, err := sql.Open("postgres", psqlInfo)
if err != nil {
log.Fatal(err)
}
defer db.Close()

// Increment the counter and get the new value
counterValue, err := incrementCounter(db)
if err != nil {
    log.Fatal(err)
}

// Generate the URL alias using base62 encoding
urlAlias := base62Encode(counterValue)
fmt.Println("Generated URL Alias:", urlAlias)
}
