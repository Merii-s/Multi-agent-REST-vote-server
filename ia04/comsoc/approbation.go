package comsoc

//SWF: retourne un décompte (count)
//SCF: retourne les alternatives préférées

//approbation: les agents votent pour autant de candidats qu'ils veulent (profils à taille variable)

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	//thresholds: nombre représentant le seuil à partir duquel les alternatives ne sont plus approuvées pour le candidat i

	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	for i := range p {
		count[p[i][0]] += 1
	}

	return count, err
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	count := make(map[Alternative]int)
	for i := range p {
		count[p[i][0]] += 1
	}

	return maxCount(count), err
}
