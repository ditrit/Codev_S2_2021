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
	"encoding/json"
	"github.com/square/certstrap/pkix"



)

func main() {
	secret := flag.String("secret","","Secret fourni par l'administrateur")
	flag.Parse()


if _, err := os.Stat("./out/client.crt"); err == nil {  // We check if the certificat exist
  connect()
  } else if os.IsNotExist(err) {
	Register(*secret)
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

func Register(secret string){
	
	// We generate the private key and the certificate request
	var key *pkix.Key
	var err error
	key, _ = pkix.CreateRSAKey(2048)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Create RSA Key error:", err)
		os.Exit(1)
	}
	keybytes,_:=key.ExportPrivate()
	private_key, err := os.Create("./out/client.key")
	if err != nil {
            panic(err)
        }
        _, err2 := private_key.WriteString(string(keybytes))
		if err2 != nil {
            panic(err)
        }
	
	stringArray := []string {"Client"}

	csr,_ := pkix.CreateCertificateSigningRequest(key, "", nil, stringArray, nil, "", "", "", "", "client")
	csrBytes, err := csr.Export()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Print certificate request error:", err)
			os.Exit(1)
		} 
	cert_request,err :=os.Create("./out/client.csr")
	if err != nil {
		panic(err)
	}
	cert_request.WriteString(string(csrBytes))

	values := map[string]string{"secret": secret, "cert_request": string(csrBytes)}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post("http://localhost:8080/cert", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalf("unable to create http request due to error %s", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if(string(body)=="Le secret fourni n'est pas le bon"){
		panic("Le secret fourni n'est pas le bon")
	}
	clientcrt, err := os.Create("./out/client.crt")
        if err != nil {
            panic(err)
        }
        clientcrt.WriteString(string(body))
		if err2 != nil {
            panic(err)
        }





	}

