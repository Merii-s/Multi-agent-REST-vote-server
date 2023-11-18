package agt

func CheckTieBreak(nb_alts int, tieBreak []int) bool {

	for i, val := range tieBreak {
		if val > nb_alts || val <= 0 {
			return true
		}
		for _, el := range tieBreak[i+1:] {
			if val == el {
				return true
			}
		}
	}
	return false
}

func Contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
