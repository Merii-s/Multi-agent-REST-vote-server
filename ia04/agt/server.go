package agt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	//rad "gitlab.utc.fr/lagruesy/ia04/demos/restagentdemo"
)

var ballotIdNumber = 0

type RestServerAgent struct {
	sync.Mutex
	id      string
	addr    string
	ballots map[string]*RestBallotAgent
}

var rules = []string{"approval", "borda", "condorcet", "copeland", "majorite_simple"}

func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{addr: addr, ballots: make(map[string]*RestBallotAgent, 0)}
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
	var date, err1 = time.Parse(time.RFC3339, req.Deadline)
	if err1 != nil || date.Before(time.Now()) || date.Equal(time.Now()) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New(" Deadline erronée"))
		return
	}
	// BadRequest : TieBreak erroné ou NombreAlts
	if CheckTieBreak(req.NbAlts, req.TieBreak) || req.NbAlts < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New(" TieBreak erroné ou NombreAlts"))
		return
	}

	// NotImplemented : Rule non implémentée
	if !Contains(rules, req.Rule) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, errors.New(fmt.Sprintf("Rule non implémentée, choisissez parmi %v", rules)))
		return
	}
	// Création du nouveau scrutin

	ballotIdNumber += 1
	ballotId := fmt.Sprintf("%s%d", "scrutin", ballotIdNumber)
	newBallot := NewRestBallotAgent(ballotId, req.Rule, req.Deadline, req.VoterIds, req.NbAlts, req.TieBreak)

	rsa.ballots[ballotId] = newBallot

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(ballotId)
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
