package agt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	rad "gitlab.utc.fr/lagruesy/ia04/demos/restagentdemo"
)

type RestVoterAgent struct {
	id    string
	url   string
	prefs []int
}

func NewRestVoterAgent(id string, url string, prefs []int) *RestVoterAgent {
	return &RestVoterAgent{id, url, prefs}
}

func (rca *RestVoterAgent) treatResponse(r *http.Response) int {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp rad.Response
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.Result
}

func (rca *RestVoterAgent) voteRequest(agentid string, voteid string, prefs []int) (err error) {
	req := RequestVote{
		Prefs: rca.prefs,
	}

	// sérialisation de la requête
	url := rca.url + "/vote"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	return
}

func (rca *RestVoterAgent) Start() {
	log.Printf("démarrage de %s", rca.id)
	err := rca.voteRequest(rca.id, "test", rca.prefs)

	if err != nil {
		log.Fatal(rca.id, "error:", err.Error())
	} else {
		log.Printf("[%s] \n", rca.id)
		log.Printf("de préférences :")
		for _, alt := range rca.prefs {
			log.Printf("%d ,", alt)
		}
	}
}
