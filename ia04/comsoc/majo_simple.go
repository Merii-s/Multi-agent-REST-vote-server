package comsoc

//SWF: retourne un décompte (count)
//SCF: retourne les alternatives préférées

func MajoritySWF(p Profile) (count Count, err error) {

	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	count = make(Count)
	for i := range p {
		count[p[i][0]] += 1
	}

	return count, err

}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	count, _ := MajoritySWF(p)
	bestAlts = maxCount(count)

	return bestAlts, nil

}
