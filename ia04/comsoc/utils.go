package comsoc

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
	return rank(alt1, prefs) > rank(alt2, prefs)
}

// renvoie les meilleures alternatives pour un décompte donné
func maxCount(count Count) (bestAlts []Alternative) {
	for alt, nb := range count {
		for _, best := range bestAlts {
			if nb > count[best] {
				bestAlts = make([]Alternative, 0)
				bestAlts = append(bestAlts, alt)
			} else if nb == count[best] {
				bestAlts = append(bestAlts, alt)
			}
		}
	}
	return
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois par préférences
func checkProfile(prefs Profile) error

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error
