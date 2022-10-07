package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func jo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "general")
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello its working")
}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Error(w, "cannot perform post request,expected GET request", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")

	if username != "" {
		userx, err := userByUsername(db, username)
		if err != nil {
			fmt.Printf("error occured: %v\n", err)
			fmt.Fprintf(w, "Error::%v", err)
			http.Error(w, "Error", http.StatusBadRequest)
			return
		}
		fmt.Println(userx)

		data, marshal_err := json.Marshal(userx)
		if marshal_err != nil {
			log.Printf("could not marshal json, err: %v", marshal_err)
			return
		}
		w.Header().Set("Content-Type", "applications/json")
		w.Write(data)
		return

	}

	// if username != "" {

	// 	if len(users) > 0 {

	// 		for _, user := range users {
	// 			if user.Username == username {
	// 				data, marshal_err := json.Marshal(user)
	// 				if marshal_err != nil {
	// 					log.Printf("could not marshal json, err: %v", marshal_err)
	// 					return
	// 				}
	// 				w.Header().Set("Content-Type", "applications/json")
	// 				w.Write(data)
	// 				return
	// 			} else {
	// 				http.Error(w, "user not found", http.StatusNotFound)
	// 				return
	// 			}

	// 		}
	// 	} else {
	// 		http.Error(w, "user not found", http.StatusNotFound)
	// 	}

	// } else {
	// 	http.Error(w, "invalid data,enter username", http.StatusBadRequest)
	// }

}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "cannot perform GET request, expected POST request", http.StatusMethodNotAllowed)
		return
	}
	//--use decoder
	decoder := json.NewDecoder(r.Body)
	var newUser User
	decode_err := decoder.Decode(&newUser)
	if decode_err != nil {
		log.Printf("error decoding: %v", decode_err)
	}
	log.Println("==========")
	log.Printf("new user created %v\n", newUser)

	user_id, err1 := addUser(db, newUser.Username, newUser.Email)
	if err1 != nil {
		fmt.Printf("error while creating user::%v\n", err1)
		http.Error(w, "name or email already exists", http.StatusBadRequest)
		return
	}
	fmt.Printf("====>>>>>new user with user_id %v created\n", user_id)
	fmt.Fprintf(w, "user added successfully")
	// // --using unmasharl
	// b, read_err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	// if read_err != nil {
	// 	fmt.Fprintf(os.Stdout, "error reading body :%v", read_err)
	// 	return
	// }

	// // Unmarshal
	// var newUser User
	// unmarshal_err := json.Unmarshal(b, &newUser)
	// if unmarshal_err != nil {
	// 	fmt.Fprintf(os.Stdout, "error unmarshilling json :%v", unmarshal_err)
	// 	return
	// }
	// fmt.Println(newUser)
}
func main() {

	http.HandleFunc("/create-user", createUserHandler)
	http.HandleFunc("/get-user", getUserHandler)
	http.HandleFunc("/store", jo)
	http.HandleFunc("retrieve", jo)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("server running on port 4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Printf("could not start server:err: %v", err)
	}
	fmt.Println("hello")
}
