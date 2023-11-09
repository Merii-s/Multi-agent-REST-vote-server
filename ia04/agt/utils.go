package agt

func checkTieBreak(nb_alts int, tieBreak []int) bool {
	checks := make(map[int]bool)
	for _, el := range tieBreak {
		if el > nb_alts || el <= 0 {
			return false
		} else if checks[el] == false {
			checks[el] = true
		} else if checks[el] == true {
			return false
		}
	}
	return true
}
