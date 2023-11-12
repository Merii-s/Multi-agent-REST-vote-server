package main

import (
	"fmt"
	"log"

	//comsoc "ia04/comsoc"
	//agt "gitlab.utc.fr/benyouzar/ai30-systeme-de-vote/-/tree/main/ia04/agt"
	agt "ai30/ia04/agt"
)

func main() {

	// lancer le serv = /newballot
	// lancer les agents en random leur prefs % option de vote du serv
	// les faire voter = /vote
	// get le res = /result
	const url1 = ":8080"
	const url2 = "http://localhost:8080"

	servAgt := agt.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go servAgt.Start()

	fmt.Scanln()
}
