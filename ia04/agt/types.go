package agt

type RequestVote struct {
	Prefs []int `json:"prefs"`
}

type ResponseVote struct {
	Error string `json:"err"`
}

type RequestNewBallot struct {
	Rule     string   `json:"rule"`
	Deadline string   `json:"deadline"`
	VoterIds []string `json:"voter_ids"`
	NbAlts   int      `json:"nb_alts"`
	TieBreak []int    `json:"tie_break"`
}

type ResponseBallot struct {
	CodeRetour string `json:"code"`
}

type Alternative int
type Profile [][]Alternative
