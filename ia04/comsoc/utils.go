package comsoc

import (
	"errors"
)

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	for i, v := range prefs {
		if v == alt {
			return i
		}
	}
	return len(prefs)
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	rank1 := rank(alt1, prefs)
	rank2 := rank(alt2, prefs)

	if rank1 == len(prefs) || rank2 == len(prefs) {
		return false
	}

	return rank1 < rank2
}

// renvoie les meilleures alternatives pour un décompte donné
func maxCount(count Count) (bestAlts []Alternative) {
	var maxCount int

	// Find the maximum count value.
	for _, nb := range count {
		if nb > maxCount {
			maxCount = nb
		}
	}

	// Collect alternatives with the maximum count.
	for alt, nb := range count {
		if nb == maxCount {
			bestAlts = append(bestAlts, alt)
		}
	}

	return
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois par préférences
func checkProfile(prefs Profile) error {
	if len(prefs) == 0 {
		return errors.New("Aucun votant pour ce scrutin")
	}
	verif := make(map[Alternative]int)
	maxLen := 0
	altUsed := make(map[Alternative]bool)

	for _, pref := range prefs {
		if len(pref) > maxLen {
			maxLen = len(pref)
		}
		for _, alt := range pref {
			if !altUsed[alt] {
				altUsed[alt] = true
			}
		}
	}

	for _, pref := range prefs { // Pour chaque votant
		for _, alt := range pref { // Pour chacune de ses prefs
			verif[alt]++
		}
		if len(verif) != maxLen {
			return errors.New("Pas même nombre d'alternatives")
		} else {
			for alt := range altUsed {
				if verif[alt] != 1 {
					return errors.New("Préférences incohérentes")
				}
			}
		}
		verif = make(map[Alternative]int) // Réinitialiser le map
	}
	return nil // Aucune erreur trouvée
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error {
	verif := make(map[Alternative]int)
	maxLen := 0
	for _, pref := range prefs {
		if len(pref) > maxLen {
			maxLen = len(pref)
		}
	}

	for _, pref := range prefs { // Pour chaque votant
		for _, alt := range pref { // Pour chacune de ses prefs
			verif[alt]++
		}
		if len(verif) != maxLen {
			return errors.New("Pas même nombre d'alternatives")
		} else {
			for _, alt := range alts {
				if verif[alt] != 1 {
					return errors.New("Préférences incohérentes")
				}
			}
		}
		verif = make(map[Alternative]int) // Réinitialiser le map
	}
	return nil // Aucune erreur trouvée
}
