// Copyright (c) 2020 Richard Youngkin. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	"os/exec"
	"strings"

)

func main() {
	flag.Parse()


if _, err := os.Stat("./out/client.crt"); err == nil {  // On vérifie si le certificat existe ou non

  connect()
  } else if os.IsNotExist(err) {
	Register()
  } 
	
}

func connect() {
	clientCertFile := flag.String("clientcert", "./out/client.crt", "Required, the name of the client's certificate file")
	clientKeyFile := flag.String("clientkey", "./out/client.key", "Required, the file name of the clients's private key file")
	srvhost := flag.String("srvhost", "localhost", "The server's host name")
	caCertFile := flag.String("cacert", "./out/ExempleCA.crt", "Required, the name of the CA that signed the server's certificate")
	var cert tls.Certificate
	var err error
	if *clientCertFile != "" && *clientKeyFile != "" {
		cert, err = tls.LoadX509KeyPair(*clientCertFile, *clientKeyFile)
		if err != nil {
			log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", *clientCertFile, *clientKeyFile)
		}
	}

	log.Printf("CAFile: %s", *caCertFile)
	caCert, err := ioutil.ReadFile(*caCertFile)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", *caCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}

	client := http.Client{Transport: t, Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s", *srvhost), bytes.NewBuffer([]byte("World")))
	if err != nil {
		log.Fatalf("unable to create http request due to error %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		switch e := err.(type) {
		case *url.Error:
			log.Fatalf("url.Error received on http request: %s", e)
		default:
			log.Fatalf("Unexpected error received: %s", err)
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("unexpected error reading response body: %s", err)
	}

	fmt.Printf("\nResponse from server: \n\tHTTP status: %s\n\tBody: %s\n", resp.Status, body)
}

func Register(){

	//secret := flag.String("secret","","Secret fourni par l'administrateur")
	// On génère la clé privé et la demande de certificat
	cmd := exec.Command("certstrap", "request-cert", "--domain", "client")
    cmd.Stdin = strings.NewReader("\n\n")
    cmd.Run()

	//_,err := http.Post("http://localhost:8080/cert","application/json",bytes.NewBuffer([]byte(*secret)))
	//if err != nil {
	//	log.Fatalf("unable to create http request due to error %s", err)
	//}
    cert_request, _ := ioutil.ReadFile("./out/client.csr")
	req,err := http.Post("http://localhost:8080/cert","application/json",bytes.NewBuffer([]byte(string(cert_request))))
	if err != nil {
		log.Fatalf("unable to create http request due to error %s", err)
	}
	body, _ := ioutil.ReadAll(req.Body)
	clientcrt, err := os.Create("./out/client.crt")
        if err != nil {
            panic(err)
        }
        _, err2 := clientcrt.WriteString(string(body))
		if err2 != nil {
            panic(err)
        }





	}


