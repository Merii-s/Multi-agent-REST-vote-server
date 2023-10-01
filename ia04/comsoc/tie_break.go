package comsoc

import "errors"

//func TieBreak([]Alternative) (Alternative, error)

func TieBreakFactory(orderedAlts []Alternative) func([]Alternative) (Alternative, error) {
	return func(alts []Alternative) (Alternative, error) {
		if len(orderedAlts) == 0 {
			return Alternative(0), errors.New("orderedAlts = nil")
		}
		fav := Alternative(0)
		for _, pref := range orderedAlts {
			for _, alt := range alts {
				if alt == pref {
					fav = alt
					break
				}
			}
		}
		return fav, nil
	}
}

func SWFFactory(swf func(Profile) (Count, error), tieBreak func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	return func(p Profile) ([]Alternative, error) {
		count, err := swf(p)
		if err != nil {
			return nil, err
		}

		best, err := tieBreak(maxCount(count))
		if err != nil {
			return nil, err
		}

		return []Alternative{best}, nil
	}
}

func SCFFactory(scf func(Profile) ([]Alternative, error), tieBreak func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	return func(p Profile) (Alternative, error) {
		bestAlts, err := scf(p)
		if err != nil {
			return Alternative(0), err
		}

		best, err := tieBreak(bestAlts)
		if err != nil {
			return Alternative(0), err
		}

		return best, nil
	}
}
