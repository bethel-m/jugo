package main

import (
	"database/sql"

	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB
var db_err error

type User struct {
	user_id  int
	Username string `json:"username"`
	Email    string `json:"email"`
}

//database connection
const (
	host     = "db"
	port     = 5432
	user     = "jugo"
	password = "jugo"
	dbname   = "jugo_db"
)

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, db_err = sql.Open("postgres", psqlInfo)
	if db_err != nil {
		fmt.Printf("could not open database: err %v\n", db_err)
		return
	}

	ping_err := db.Ping()
	if ping_err != nil {
		fmt.Printf("ping error :%v\n", ping_err)
		return
	}
	fmt.Println("connection_successful")

}

func userByUsername(db *sql.DB, name string) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE name=$1", name)
	if err := row.Scan(&user.user_id, &user.Username, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("userByUsername: %v, no such user", name)
		}
		return user, fmt.Errorf("userByUsername %v:%v", name, err)
	}
	return user, nil
}

func addUser(db *sql.DB, name string, email string) (int, error) {
	var user_id int
	// _, err := db.Exec("INSERT INTO users (name,email) VALUES($1,$2)", name, email)
	// if err != nil {
	// 	return 0, fmt.Errorf("addUser: %v", err)
	// }

	// return 1, nil

	err := db.QueryRow("INSERT INTO users(name,email) VALUES($1,$2) RETURNING user_id",
		name, email).Scan(&user_id)
	if err != nil {
		return 0, err
	}
	return user_id, err
}
