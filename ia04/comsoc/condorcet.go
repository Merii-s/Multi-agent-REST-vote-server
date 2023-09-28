package comsoc

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	return bestAlts, nil
}
