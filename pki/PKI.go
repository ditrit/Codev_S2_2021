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
    "strings"
    "unicode"
)
type User struct {
	Login   string  `json:"login"`
	Password int64 `json:"password"`
}
type Secret struct {
    Secret string `json:"secret"`
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
            write_secret("secrets.txt",secret,true)


		} else {
			w.Header().Set("Contotto)ent-Type", "application/json")
    		w.WriteHeader(http.StatusCreated)
    		w.Write([]byte(`Mauvais identifiants`))
		}


    
}

func cert(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received %s request for host %s from IP address %s and X-FORWARDED-FOR %s",
			r.Method, r.Host, r.RemoteAddr, r.Header.Get("X-FORWARDED-FOR"))
    var s Secret
		body, err := ioutil.ReadAll(r.Body)
        fmt.Print("in")
        fmt.Print(s.Secret,"oui")
        fmt.Print(string(body))
		if err != nil {
			body = []byte(fmt.Sprintf("error reading request body: %s", err))
		}

        if(stringInSlice(secret,get_secret("secrets.txt"))){
            write_secret("secrets.txt",secret,false)
            fmt.Print("secret supprimé")
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
    r := mux.NewRouter()
    r.HandleFunc("/", get).Methods(http.MethodGet)
    r.HandleFunc("/login", post).Methods(http.MethodPost)
    r.HandleFunc("/cert", cert).Methods(http.MethodPost)
    r.HandleFunc("/", notFound)
    log.Fatal(http.ListenAndServe(":8080", r))
}

func get_secret(secrets string) []string{ //fonction permettant de transformer le fichier text en list de secrets
	data, err := ioutil.ReadFile(secrets)
	if err != nil {
	fmt.Println("File reading error", err)
}
f := func(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
}
var secrets_list=strings.FieldsFunc(string(data), f)
return secrets_list
}

func write_secret(secrets string,code string,a bool){ //fonction permettant de transformer le fichier text en list de secrets
    data, _ := ioutil.ReadFile(secrets)
    if a {
        var string_secret = string(data)
        string_secret += ", "+code
        ioutil.WriteFile(secrets, []byte(string_secret), 0644)

    }else {
        var liste_secret= get_secret(secrets)
        if(stringInSlice(code,liste_secret)) {
            liste_secret=del_el(code,liste_secret)
            ioutil.WriteFile(secrets, []byte(strings.Join(liste_secret, " ")), 0644)
        }
        

    }

	
}
func stringInSlice(a string, list []string) bool { //Vérifie la présence d'un élement dans une liste
    for _, b := range list {
        if b == a {
            return true}
    }
    return false
}

func del_el(a string, list []string ) []string { // Permet de supprimer un élement d'une liste
	var c=-1
	for i,n:= range list {
		if n == a {
			c=i
		}
	}
	if c!=-1 {
		return append(list[:c],list[c+1:]...)
	}

	return list
}

