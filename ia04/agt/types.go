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
	VoterIds []string `json:"voter-ids"`
	NbAlts   int      `json:"alts"`
	TieBreak []int    `json:"tie-break"`
}

type ResponseBallot struct {
	CodeRetour string `json:"code"`
}
