package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v9"
)

var client *redis.Client

//unimplemented function
func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "general")
}

// this function gets the user from the database
// this endpoint recieves a username parameter from the client
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	//allows only get request on this endpoint
	if r.Method != "GET" {
		http.Error(w, "cannot perform request,expected GET request", http.StatusMethodNotAllowed)
	}

	//get username from the parameters sent in the request
	username := r.URL.Query().Get("username")

	//query the database if the username parameter is not nill
	if username != "" {
		userx, err := userByUsername(db, username)
		if err != nil {
			fmt.Printf("error occured: %v\n", err)
			http.Error(w, "user does not exist", http.StatusBadRequest)
			return
		}
		fmt.Println(userx)

		//convert user struct to a json object and send it to the client
		data, marshal_err := json.Marshal(userx)
		if marshal_err != nil {
			log.Printf("could not marshal json, err: %v", marshal_err)
			return
		}
		w.Header().Set("Content-Type", "applications/json")
		w.Write(data)
		return

	}

}

// this function handles the /store endpoint
// this endpoint recieves a file(multipart/form) from the client,and a username parameter
func storeHandler(w http.ResponseWriter, r *http.Request) {
	// only POST method is allowed on this endpoint
	if r.Method != http.MethodPost {
		http.Error(w, "post method only,allowed", http.StatusBadRequest)
		return
	}

	//check if user exists in database
	//if the user does not exist an error is sent to the client
	//if the user exists the function continues
	username := r.URL.Query().Get("username")
	_, err := userByUsername(db, username)
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}

	//parse the file sent from the client
	multipart_err := r.ParseMultipartForm(32 << 20)
	if multipart_err != nil {
		fmt.Fprintf(os.Stdout, "parsemultiform error: %v", multipart_err)
		return
	}
	//reads the contents of the file from the user into file,
	//and metadata about the file into handler
	// any errors are stored in form_err
	file, handler, form_err := r.FormFile("file")
	if form_err != nil {
		fmt.Fprintf(os.Stdout, "formfile error::%v", form_err)
		http.Error(w, "could not parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// create a directory to store the files the client sends
	// if the directory exists skip the creation and use it
	// send a server error to the client in case of error and return
	mkdir_err := os.MkdirAll("./uploads", os.ModePerm)
	if mkdir_err != nil {
		fmt.Fprintf(os.Stdout, "could not make dir uploads ::%v", mkdir_err)
		http.Error(w, mkdir_err.Error(), http.StatusInternalServerError)
		return
	}

	// create a file in the uploads directory with this naming convention "username_filename"
	// using the username and the filename frpm the request, copy the contents of the sent
	// file to this file and save it
	dst, create_err := os.Create(fmt.Sprintf("./uploads/%s_%s", username, handler.Filename))
	if create_err != nil {
		fmt.Fprintf(os.Stdout, "could not create file ::%v", create_err)
		http.Error(w, create_err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()
	//copy the sent file to the created file
	//if and error occurs while copying the file ,send server error to the client and return
	_, copy_err := io.Copy(dst, file)
	if copy_err != nil {
		fmt.Fprintf(os.Stdout, "error copying file to destiantion::%v", copy_err)
		http.Error(w, copy_err.Error(), http.StatusInternalServerError)
		return
	}

	// add the newfile name to the redis queue using the add_task_to_queue function
	task_name := username + "_" + handler.Filename
	numberOfTask, err_in_queue := add_task_to_queue(client, task_name)
	if err_in_queue != nil {
		fmt.Fprintf(os.Stdout, "error adding task to queue:::%v", err_in_queue)
		return
	}

	fmt.Println(task_name)
	fmt.Printf("%v tasks in queue\n", numberOfTask)
	fmt.Fprintf(w, "%v moved for processing\n", handler.Filename)
	fmt.Println("file has been passed for processing")
}

// this endpoint creates a user and stores the user instance in the database
// it takes json request of username and email and creates a user
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	//only post method is allowed
	if r.Method != "POST" {
		http.Error(w, "cannot perform GET request, expected POST request", http.StatusMethodNotAllowed)
		return
	}

	// unmarshal the json from the request body into the user struct
	// using a decoder,if an error occurs send a server error and return
	decoder := json.NewDecoder(r.Body)
	var newUser User
	decode_err := decoder.Decode(&newUser)
	if decode_err != nil {
		http.Error(w, "could not read json::", http.StatusInternalServerError)
		log.Printf("error decoding: %v", decode_err)
		return
	}
	log.Println("==========")
	log.Printf("new user instance created %v\n", newUser)

	// using the addUser funcion,persist the new user instance
	// to the database,if an error occurs send a server error to the client
	// and return
	user_id, err1 := addUser(db, newUser.Username, newUser.Email)
	if err1 != nil {
		fmt.Printf("error while creating user::%v\n", err1)
		http.Error(w, "name or email already exists", http.StatusInternalServerError)
		return
	}

	//confirmation message when the user gets added successfully
	fmt.Printf("====>>>>>new user with user_id %v created\n", user_id)
	fmt.Fprintf(w, "user added successfully")

}

func main() {
	//endpoints for the program
	http.HandleFunc("/create-user", createUserHandler)
	http.HandleFunc("/get-user", getUserHandler)
	http.HandleFunc("/store", storeHandler)
	http.HandleFunc("/retrieve", retrieveHandler)

	fmt.Println("server running on port 4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Printf("could not start server:err: %v\n", err)
	}
	fmt.Println("hello")
}
