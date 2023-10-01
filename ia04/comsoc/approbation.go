package comsoc

import "errors"

//SWF: retourne un décompte (count)
//SCF: retourne les alternatives préférées

//approbation: les agents votent pour autant de candidats qu'ils veulent (profils à taille variable)

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	//thresholds: nombre représentant le seuil à partir duquel les alternatives ne sont plus approuvées pour le candidat i

	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	if len(thresholds) != len(p) {
		return nil, errors.New("len(tresholds) != len(p)")
	}

	count = make(Count)

	//initialisation pour que les candidats avec zéro voix apparaissent dans le décompte
	for _, pref := range p {
		for _, alt := range pref {
			count[alt] = 0
		}
	}

	for i, pref := range p {
		for j := 0; j < thresholds[i]; j++ {
			count[pref[j]] += 1
		}
	}

	return count, err

}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	count, _ := ApprovalSWF(p, thresholds)
	bestAlts = maxCount(count)

	return bestAlts, err
}
