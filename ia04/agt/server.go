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
	// const avec les règles ?

	// Création du nouveau scrutin

	ballotIdNumber += 1
	ballotId := fmt.Sprintf("%s%d", "scrutin", ballotIdNumber)
	newBallot := NewRestBallotAgent(ballotId, req.Rule, req.Deadline, req.VoterIds, req.NbAlts, req.TieBreak)

	rsa.ballots[ballotId] = newBallot

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(ballotId)
	w.Write(serial)
}

func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	// lock le serveur
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	var voteReq VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&voteReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Vérification si le scrutin existe
	if _, ok := rsa.ballots[voteReq.BallotID]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Scrutin non trouvé")
		return
	}

	// Vérification de la deadline du scrutin
	ballot, ok := rsa.ballots[voteReq.BallotID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Scrutin non trouvé")
		return
	}
	deadline, _ := time.Parse(time.RFC3339, ballot.deadline)
	if time.Now().After(deadline) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "La deadline est passée")
		return
	}

	// Vérification si l'agent a déjà voté
	if rsa.hasAgentVoted(voteReq.AgentID, voteReq.BallotID) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Vote de cet agent déjà effectué")
		return
	}

	// NotImplemented : Vérification d'options facultatives ici (par exemple, seuil d'acceptation en approval)
	if CheckOptions(voteReq.Options) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New(" Options erronées"))
		return
	}

	// Effectuer le vote
	ballot.Vote(voteReq.AgentID, voteReq.Prefs, voteReq.Options)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Vote enregistré")
}

// hasAgentVoted vérifie si l'agent a déjà voté pour ce scrutin
func (rsa *RestServerAgent) hasAgentVoted(agentID, ballotID string) bool {
	// Vérifie si l'agent est dans la liste des votants pour ce scrutin
	ballot, _ := rsa.ballots[ballotID]
	for _, voterID := range ballot.voter_ids {
		if voterID == agentID {
			return true
		}
	}
	return false
}

// Cette méthode doit être implémentée dans RestBallotAgent pour permettre le vote
func (rba *RestBallotAgent) Vote(agentID string, pref, options []int) {
	// Si l'agent est autorisé à voter, enregistrement du vote en ajoutant l'agent à la liste des votants
	rba.voter_ids = append(rba.voter_ids, agentID)
	fmt.Println("Vote enregistré pour l'agent", agentID)
}

// CheckOptions vérifie si les options sont valides
func CheckOptions(options []int) bool {
	// Implémentez votre logique pour vérifier si les options sont valides
	// Retourne true si les options sont valides, sinon false
	return false
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.doNewBallot)
	mux.HandleFunc("/vote", rsa.doVote) // Ajout de l'endpoint /vote

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
