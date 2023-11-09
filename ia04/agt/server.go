package agt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	//"ai30-systeme-de-vote/ia04/comsoc"

	rad "gitlab.utc.fr/lagruesy/ia04/demos/restagentdemo"
)

var ballotIdNumber = 0

type RestServerAgent struct {
	sync.Mutex
	id      string
	addr    string
	ballots map[string]*RestBallotAgent
}

func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{addr: addr}
}

// Test de la méthode
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestServerAgent) decodeRequest(r *http.Request) (req RequestNewBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestServerAgent) doNewBallot(w http.ResponseWriter, r *http.Request) {
	// lock le serv
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	// BadRequest : Deadline passée
	if req.Deadline < time.Now().String() || req.NbAlts < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	// BadRequest : TieBreak erroné
	if checkTieBreak(req.NbAlts, req.TieBreak) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// BadRequest : voter_ids inconnus ?

	// NotImplemented : Rule non implémentée
	// const avec les règles ?

	// Création du nouveau scrutin
	var resp rad.Response

	ballotIdNumber += 1
	ballotId := fmt.Sprintf("%s%d", "scrutin", ballotIdNumber)
	newBallot := NewRestBallotAgent(ballotId, req.Rule, req.Deadline, req.VoterIds, req.NbAlts, req.TieBreak)

	rsa.ballots[ballotId] = newBallot

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.doNewBallot)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
