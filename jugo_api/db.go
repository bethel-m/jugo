package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fmt"

	"github.com/joho/godotenv"
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

// init function to initialize database
func init() {
	// check if there exist a .env file
	// if there is,load the environmental variables
	// and continue
	current_directory, dir_err := os.Getwd()
	if dir_err != nil {
		log.Fatalf("could not get current working directory>>%v", dir_err)
	}
	parent_directory := filepath.Dir(current_directory)
	env_path := filepath.Join(parent_directory, ".env")
	if env_exists, file_err := fileExists(env_path); env_exists {
		if file_err == nil {
			err := godotenv.Load(env_path)
			if err != nil {
				log.Fatalf("Error loading .env file>>%v", err)
			}
		} else {
			log.Fatalf("error in config file>>:%v", file_err)
		}

	}

	// get database variables from .env file
	host := os.Getenv("DB_HOST")
	port, port_err := strconv.Atoi(os.Getenv("DB_PORT"))
	if port_err != nil {
		log.Fatalf("error converting port to int::%v", port_err)
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

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
		log.Fatalf("database ping error :%v\n", ping_err)
	}
	fmt.Println("connection_successful")

}

// this function takes a database struct and a username
// queries the database and returns the user or and error
func userByUsername(db *sql.DB, name string) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE name=$1", name)
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

	err := db.QueryRow("INSERT INTO users(name,email) VALUES($1,$2) RETURNING user_id",
		name, email).Scan(&user_id)
	if err != nil {
		return 0, err
	}
	return user_id, err
}

// this function takes a database struct,and a username
// it queries the database for the file path of the username
// it return the filepath of the username or error otherwise
func retrieveFilePath(db *sql.DB, name string) (string, error) {
	var filepath string
	row := db.QueryRow("SELECT file_path FROM users WHERE name=$1", name)
	if err := row.Scan(&filepath); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("retrieveFilePath: %v ,no such user", name)
		}
		return "", fmt.Errorf("retrieveFilePath:: %v:%v", name, err)
	}
	return filepath, nil
}

// this function checks if a file exists or not
// it returns a boolean and an error,the boolean value specifies if the
// file exists or not ,while the error value specifies any error
// if the file exists it returns true and nil,if the file does not exist
// it returns false and nil,if the file exists but an error occured it returns
// true and an error
func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err
	}
	return true, nil
}
