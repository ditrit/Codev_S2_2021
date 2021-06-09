# Codev_S2_2021

<p>
Gandalf is an open-source solution that automates the construction of the (complex) tooling needed to develop products in DevOps mode.
Gandalf is distributed and dynamic (the nature and number of its components can change at any time) and its communications must be secure.</p>

<p>The project consists of implementing the dynamic management of certificates
used to secure communications between Gandalf components.
It is a matter of fully automating (without human intervention) the process of
certificate request/provision process, including the validation stage of the request.</p>

<h2> Our system is composed of four independent blocks: </h2>

<img src="https://user-images.githubusercontent.com/83370247/121272276-d5710b00-c8c5-11eb-9431-8b42bf91d241.png" >

<h2> How to use it: </h2>

```sh
git clone https://github.com/ditrit/Codev_S2_2021.git
cd Codev_S2_2021
```
<h3>Start all necessary components:</h4>

```sh
cd ./pki
go run pki.go

cd ../serveur
go run serveur.go
```

<h3> Get a secret: </h3>

```sh
curl -X POST http://localhost:8080/login -d "{""login"":""admin"",""password"":123}"
```
<h3
<h3> Initialize and get the certificat of a composant: </h3>

```sh

cd ./client
go run client.go --secret "secret"

```

<h3> Connect to the Gandalf serveur </h3>

```sh

cd ./client
go run client.go

```
