package agt

import (
	comsoc "ai30/ia04/comsoc"
	"fmt"
	"log"
)

type RestBallotAgent struct {
	id          string
	rule        string
	deadline    string
	voter_ids   map[string]bool // true si a voté
	alts_number int
	tie_break   []int
	profile     [][]int
	options     []int
}

func NewRestBallotAgent(id string, rule string, deadline string, voter_ids map[string]bool, alts_number int, tie_break []int, profile [][]int) *RestBallotAgent {
	return &RestBallotAgent{
		id:          id,
		rule:        rule,
		deadline:    deadline,
		voter_ids:   voter_ids,
		alts_number: alts_number,
		tie_break:   tie_break,
		profile:     profile,
	}
}

// Check deadline
// Check nb altertive
// Check si a pas déjà voté

func (rca *RestBallotAgent) Start() {
	log.Printf("Nouveau scrutin : %s", rca.id)
}

// Enregistrement du vote
func (rba *RestBallotAgent) Vote(agentID string, pref, options []int) {
	// Ajout des préférences au profil (on ne prend pas en compte l'ordre des agt_id)
	//rba.profile[rba.votedNumber-1] = pref
	rba.profile = append(rba.profile, pref)

	// Marquer l'agent comme ayant voté
	rba.voter_ids[agentID] = true
	fmt.Println("Vote enregistré pour l'agent", agentID)
	// fmt.Println("Profil après enregistrement du vote de l'agent", rba.profile)
}

// Obtention du gagnant
func (rba *RestBallotAgent) GetWinner() (winner comsoc.Alternative, ranking []comsoc.Alternative) {
	// Conversion du profil en [][]Alternative
	profileAlt := convertToAlternative(rba.profile)
	tieBreak := comsoc.TieBreakFactory(convertToAlternativeTieBreak(rba.tie_break))

	var swf func(p comsoc.Profile) ([]comsoc.Alternative, error)
	var scf func(p comsoc.Profile) (comsoc.Alternative, error)

	switch rba.rule {
	case "majority":
		scf = comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak)
		swf = comsoc.SWFFactory(comsoc.MajoritySWF, tieBreak)

	case "borda":
		scf = comsoc.SCFFactory(comsoc.BordaSCF, tieBreak)
		swf = comsoc.SWFFactory(comsoc.BordaSWF, tieBreak)

	case "condorcet":
		scf = comsoc.SCFFactory(comsoc.CondorcetWinner, tieBreak)
		swf = nil

	case "copeland":
		scf = comsoc.SCFFactory(comsoc.CopelandWinner, tieBreak)
		swf = nil

	case "approval":
		ranking := make([]comsoc.Alternative, 0)
		count, _ := comsoc.ApprovalSWF(profileAlt, rba.options)
		for len(count) > 0 {
			best_alts := comsoc.MaxCount(count)
			best, _ := tieBreak(best_alts)
			ranking = append(ranking, best)
			delete(count, best)
		}
		return ranking[0], ranking
	}

	winner, _ = scf(profileAlt)
	ranking, _ = swf(profileAlt)

	return winner, ranking
}
