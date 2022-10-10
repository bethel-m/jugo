package main

import (
	"database/sql"
	"time"

	"fmt"

	_ "github.com/lib/pq"
)

//database struct ,passed to different endpoints when required
var db *sql.DB

// database error passed to the different function or enpoints when required
var db_err error

// user struct to be stored and retrieved
type User struct {
	user_id      int
	Username     string `json:"username"`
	Email        string `json:"email"`
	Last_changed time.Time
	File_path    string
}

//database connection variables
const (
	host     = "localhost"
	port     = 5432
	user     = "jugo"
	password = "jugo"
	dbname   = "jugo_db"
)

// init function to initialize database
func init() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	//connecting to the database
	db, db_err = sql.Open("postgres", psqlInfo)
	if db_err != nil {
		fmt.Printf("could not open database: err %v", db_err)
	}
	ping_err := db.Ping()
	if ping_err != nil {
		fmt.Printf("ping error :%v", ping_err)
	}
	fmt.Println("connection_successful")

}

// this function takes a database struct and a username
// queries the database and returns the user or and error
func userByUsername(db *sql.DB, name string) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM users_with_filespath WHERE name=$1", name)
	if err := row.Scan(&user.user_id, &user.Username, &user.Email, &user.Last_changed, &user.File_path); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("userByUsername: %v, no such user", name)
		}
		return user, fmt.Errorf("userByUsername %v:%v", name, err)
	}
	return user, nil
}

// this function takes a database struct, a username and an email
// it creates inserts /persist the user to the database
// it returns the users id if successful,or an error otherwise
func addUser(db *sql.DB, name string, email string) (int, error) {
	var user_id int

	err := db.QueryRow("INSERT INTO users_with_filespath(name,email) VALUES($1,$2) RETURNING user_id",
		name, email).Scan(&user_id)
	if err != nil {
		return 0, err
	}
	return user_id, err
}
