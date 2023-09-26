package comsoc

//SWF: retourne un décompte (count)
//SCF: retourne les alternatives préférées

//Borda: 0 pour le dernier, m-1 (m nombre de candidats) pour le premier

func BordaSWF(p Profile) (count Count, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	count = make(Count)

	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			count[p[i][j]] += len(p[i]) - j - 1
		}
	}

	return count, err
}

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	count := make(Count)
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			count[p[i][j]] += len(p[i]) - 1 - j
		}
	}

	return maxCount(count), err
}
