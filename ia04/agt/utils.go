package agt

func CheckAlternativeConsistency(nb_alts int, tieBreak []int) bool {

	verif := make(map[int]int)

	// Alternative hors range
	for _, val := range tieBreak {
		if val > nb_alts || val <= 0 {
			return true
		}
	}

	for _, alt := range tieBreak { // Pour chacune de ses prefs
		verif[alt]++
	}
	// Pas le même nombre d'alternatives
	if len(verif) != nb_alts {
		return true
	} else {
		// Chaque élément présent est bien entre 1 et nb_alts et apparaît une seule fois
		for i := 1; i <= nb_alts; i++ {
			if verif[i] != 1 {
				return false
			}
		}
	}
	return false
}

// return true if value is in arr
func Contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
