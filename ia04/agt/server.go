package agt

import (
	comsoc "ai30/ia04/comsoc"
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

var rules = []string{"approval", "borda", "condorcet", "copeland", "majority"}

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
	if CheckAlternativeConsistency(req.NbAlts, req.TieBreak) || req.NbAlts < 1 {
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

	// Récupération des IDs des votants dans une map
	voterIds := make(map[string]bool)
	for _, voterId := range req.VoterIds {
		voterIds[voterId] = false
	}

	// Initialisation des préférences avec une structure Profile vide pour chaque votant
	profile := make([][]int, 0)
	for i := range profile {
		profile[i] = make([]int, req.NbAlts)
	}

	newBallot := NewRestBallotAgent(ballotId, req.Rule, req.Deadline, voterIds, req.NbAlts, req.TieBreak, profile)

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

	// Vérification de l'existence du scrutin
	ballot, ok := rsa.ballots[voteReq.BallotID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Scrutin non trouvé")
		return
	}

	// Vérification de la deadline du scrutin
	deadline, _ := time.Parse(time.RFC3339, ballot.deadline)
	if time.Now().After(deadline) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "La deadline est passée")
		return
	}

	// Vérification si l'agent a le droit de voter
	if _, exists := ballot.voter_ids[voteReq.AgentID]; !exists {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Agent non autorisé à voter pour ce scrutin")
		return
	}

	// Vérification si agent a déjà voté
	if ballot.voter_ids[voteReq.AgentID] {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Vote de cet agent déjà effectué")
		return
	}

	// Vérification des préférences du votant
	if CheckAlternativeConsistency(rsa.ballots[voteReq.BallotID].alts_number, voteReq.Prefs) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Préférences erronnées")
		return
	}

	// Vérification des options
	if ballot.rule == "approval" {
		if len(voteReq.Options) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Options non renseignées pour la rule 'approval'")
			return
		}
		rsa.ballots[voteReq.BallotID].options = append(rsa.ballots[voteReq.BallotID].options, voteReq.Options[0])
	} else {
		if len(voteReq.Options) > 0 {
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprint(w, "Options spécifiées pour une rule autre que 'approval'")
			return
		}
	}

	// Effectuer le vote
	ballot.Vote(voteReq.AgentID, voteReq.Prefs, voteReq.Options)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Vote enregistré")
}

func (rsa *RestServerAgent) doResult(w http.ResponseWriter, r *http.Request) {
	// lock le serveur
	rsa.Lock()
	defer rsa.Unlock()

	// Vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	var resultReq ResultRequest
	if err := json.NewDecoder(r.Body).Decode(&resultReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Vérification de l'existence du scrutin
	ballot, ok := rsa.ballots[resultReq.BallotID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Scrutin non trouvé")
		return
	}

	// Vérification de la deadline du scrutin
	deadline, _ := time.Parse(time.RFC3339, ballot.deadline)
	if time.Now().Before(deadline) {
		w.WriteHeader(http.StatusTooEarly)
		fmt.Fprint(w, "La deadline n'est pas encore passée")
		return
	}

	winner, ranking := ballot.GetWinner()

	if winner == -1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Gagnant non trouvé")
		return
	}

	// Convert []Alternative to []int
	var intRanking []int
	for _, alt := range ranking {
		intRanking = append(intRanking, int(alt))
	}

	response := map[string]interface{}{
		"winner":  int(winner),
		"ranking": intRanking,
	}

	serial, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(serial)
}

// Conversion de [][]int en [][]Alternative
func convertToAlternative(profile [][]int) comsoc.Profile {
	result := make(comsoc.Profile, len(profile))
	for i := range profile {
		result[i] = make([]comsoc.Alternative, len(profile[i]))
		for j := range profile[i] {
			result[i][j] = comsoc.Alternative(profile[i][j])
		}
	}
	// fmt.Println("Profil converti en [][]Alternative : ", result)
	return result
}

func convertToAlternativeTieBreak(profile []int) []comsoc.Alternative {
	result := make([]comsoc.Alternative, 0)
	for _, alt := range profile {
		result = append(result, comsoc.Alternative(alt))
	}
	return result
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.doNewBallot)
	mux.HandleFunc("/vote", rsa.doVote)
	mux.HandleFunc("/result", rsa.doResult)

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
