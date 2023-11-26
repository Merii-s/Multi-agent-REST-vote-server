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
	votedNumber int
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
	rba.votedNumber += 1
	//rba.profile[rba.votedNumber-1] = pref
	rba.profile = append(rba.profile, pref)

	// Marquer l'agent comme ayant voté
	rba.voter_ids[agentID] = true
	fmt.Println("Vote enregistré pour l'agent", agentID)
	// fmt.Println("Profil après enregistrement du vote de l'agent", rba.profile)
}

// Obtention du gagnant
func (rba *RestBallotAgent) GetWinner() comsoc.Alternative {
	// Conversion du profil en [][]Alternative
	profileAlt := convertToAlternative(rba.profile)

	switch rba.rule {
	case "majority":
		tieBreak := comsoc.TieBreakFactory(convertToAlternativeTieBreak(rba.tie_break))
		SCFMajo := comsoc.SCFFactory(comsoc.MajoritySCF, tieBreak)
		bestAlt, err := SCFMajo(profileAlt)
		if err != nil {
			return -1
		}
		return bestAlt
	case "borda":
		bestAlts, err := comsoc.BordaSCF(profileAlt)
		if err != nil || len(bestAlts) == 0 {
			return -1
		}
		return bestAlts[0]
	case "approval":
		bestAlts, err := comsoc.ApprovalSCF(profileAlt, rba.tie_break)
		if err != nil || len(bestAlts) == 0 {
			return -1
		}
		return bestAlts[0]
	}

	return -1
}
