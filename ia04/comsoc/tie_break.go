package comsoc

import (
	"errors"
	"sort"
)

//func TieBreak([]Alternative) (Alternative, error)

func TieBreakFactory(orderedAlts []Alternative) func([]Alternative) (Alternative, error) {
	return func(alts []Alternative) (Alternative, error) {
		if len(orderedAlts) == 0 {
			return Alternative(0), errors.New("orderedAlts = nil")
		}
		fav := Alternative(0)
	Loop:
		for _, pref := range orderedAlts {
			for _, alt := range alts {
				if alt == pref {
					fav = alt
					break Loop
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

		// Create a slice of keys from the map
		var rankedAlts []Alternative
		for key := range count {
			rankedAlts = append(rankedAlts, key)
		}

		// Sort the keys based on map values in descending order
		sort.Slice(rankedAlts, func(i, j int) bool {
			return count[rankedAlts[i]] > count[rankedAlts[j]]
		})

		// Create a slice of integers in descending order based on map values
		var result []Alternative
		for _, key := range rankedAlts {
			result = append(result, key)
		}

		return rankedAlts, nil
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
