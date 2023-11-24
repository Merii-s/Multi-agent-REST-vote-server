package agt

import (
	"log"
)

type RestBallotAgent struct {
	id          string
	rule        string
	deadline    string
	voter_ids   map[string]bool // true si a voté
	alts_number int
	tie_break   []int
}

func NewRestBallotAgent(id string, rule string, deadline string, voter_ids map[string]bool, alts_number int, tie_break []int) *RestBallotAgent {
	return &RestBallotAgent{id, rule, deadline, voter_ids, alts_number, tie_break}
}

// Check deadline
// Check nb altertive
// Check si a pas déjà voté

func (rca *RestBallotAgent) Start() {
	log.Printf("Nouveau scrutin : %s", rca.id)
}
