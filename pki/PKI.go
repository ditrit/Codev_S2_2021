package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
	"encoding/json"
	"fmt"
    "github.com/google/uuid"
    "io/ioutil"
    "os"
    "os/exec"


)
type User struct {
	Login   string  `json:"login"`
	Password int64 `json:"password"`
}
var secret=uuid.New().String()

func get(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	var u User
        if r.Body == nil {
            http.Error(w, "Please send a request body", 400)
            return
        }
        err := json.NewDecoder(r.Body).Decode(&u)
        if err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
		if u.Login=="admin" && u.Password==123{
			w.Header().Set("Contotto)ent-Type", "application/json")
    		w.WriteHeader(http.StatusCreated)
    		secret=uuid.New().String()
            w.Write([]byte(secret))
		} else {
			w.Header().Set("Contotto)ent-Type", "application/json")
    		w.WriteHeader(http.StatusCreated)
    		w.Write([]byte(`Mauvais identifiants`))
		}


    
}

func cert(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received %s request for host %s from IP address %s and X-FORWARDED-FOR %s",
			r.Method, r.Host, r.RemoteAddr, r.Header.Get("X-FORWARDED-FOR"))
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			body = []byte(fmt.Sprintf("error reading request body: %s", err))
		}
        // recreate cliente csr
        clientcsr, err := os.Create("./out/client.csr")
        if err != nil {
            panic(err)
        }
        //Ne pas arreter le serveur
        _, err2 := clientcsr.WriteString(string(body))
        if err2 != nil {
            panic(err)
        }
        clientcsr.Close()      
        // On signe le certificat
        cmd2 := exec.Command("certstrap", "sign", "client", "--CA","ExempleCA")
        cmd2.Run()


        // On le renvoie
        signed_cert, _ := ioutil.ReadFile("./out/client.crt")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write(signed_cert)
}

func notFound(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(`{"message": "not found"}`))
}

func main() {
	var secret int64
    fmt.Println(secret)
    r := mux.NewRouter()
    r.HandleFunc("/", get).Methods(http.MethodGet)
    r.HandleFunc("/login", post).Methods(http.MethodPost)
    r.HandleFunc("/cert", cert).Methods(http.MethodPost)
    r.HandleFunc("/", notFound)
    log.Fatal(http.ListenAndServe(":8080", r))
}