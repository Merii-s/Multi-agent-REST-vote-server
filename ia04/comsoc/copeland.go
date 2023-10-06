package comsoc

func CopelandtWinner(p Profile) (bestAlts []Alternative, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}
	nbalt := len(p[0])
	// Initialiser un tableau pour compter les victoires de chaque candidat
	wins := make(Count)
	// Initialiser un tableau pour compter les victoires durant le 1v1
	win1v1 := make(Count)

	//Pour chaque 1v1 possible
	for i := 1; i <= nbalt; i++ {
		for j := 1; j <= nbalt; j++ {
			if i != j {
				// On fait le 1v1 sur le profil
				for _, prefs := range p {
					if isPref(Alternative(i), Alternative(j), prefs) {
						win1v1[Alternative(i)]++
					} else {
						win1v1[Alternative(j)]++
					}
				}
				if win1v1[Alternative(i)] > win1v1[Alternative(j)] {
					wins[Alternative(i)]++
				} else {
					wins[Alternative(i)]--
				}
				// Réinitialise le nombre de victoire pour le prochain 1v1
				win1v1 = make(Count)
			}
		}
	}

	// Trouver les candidats ayant remporté le plus de comparaisons
	bestAlts = append(bestAlts, Alternative(-len(p[0])))
	for alt, nbwin := range wins {
		if nbwin == wins[Alternative(bestAlts[0])] && bestAlts[0] != Alternative(-len(p[0])) {
			bestAlts = append(bestAlts, alt)
		} else if nbwin > wins[Alternative(bestAlts[0])] || bestAlts[0] == Alternative(-len(p[0])) {
			bestAlts = make([]Alternative, 0)
			bestAlts = append(bestAlts, alt)
		}
	}

	return bestAlts, nil

}
