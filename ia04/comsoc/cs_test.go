// version 2.0.0

package comsoc

import (
	"fmt"
	"testing"
)

func TestBordaSWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	prefs2 := [][]Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 3, 2, 4},
		{1, 3, 2, 4},
		{1, 3, 2, 4},
		{1, 3, 2, 4},
		{4, 2, 1, 3},
		{4, 2, 1, 3},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{4, 3, 2, 1},
		{4, 3, 2, 1},
	}

	res, _ := BordaSWF(prefs)
	res2, _ := BordaSWF(prefs2)

	if res[1] != 4 {
		t.Errorf("error, result for 1 should be 4, %d computed", res[1])
	}
	if res[2] != 3 {
		t.Errorf("error, result for 2 should be 3, %d computed", res[2])
	}
	if res[3] != 2 {
		t.Errorf("error, result for 3 should be 2, %d computed", res[3])
	}
	if res2[1] != 37 {
		t.Errorf("error, result for 1 should be 37, %d computed", res2[1])
	}
	if res2[2] != 48 {
		t.Errorf("error, result for 2 should be 48, %d computed", res2[2])
	}
	if res2[3] != 47 {
		t.Errorf("error, result for 3 should be 47, %d computed", res2[3])
	}
	if res2[4] != 30 {
		t.Errorf("error, result for 4 should be 30, %d computed", res2[4])
	}
}

func TestBordaSCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := BordaSCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestMajoritySWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	prefs2 := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{2, 3, 1},
		{2, 3, 1},
		{2, 3, 1},
		{2, 3, 1},
		{2, 3, 1},
		{2, 3, 1},
		{3, 2, 1},
		{3, 2, 1},
		{3, 2, 1},
		{3, 2, 1},
		{3, 2, 1},
	}

	res, _ := MajoritySWF(prefs)
	res2, _ := MajoritySWF(prefs2)

	// fmt.Print(res)
	// fmt.Print(res2)

	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 0 {
		t.Errorf("error, result for 2 should be 0, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
	if res2[1] != 10 {
		t.Errorf("error, result for 1 should be 10, %d computed", res2[1])
	}
	if res2[2] != 6 {
		t.Errorf("error, result for 2 should be 6, %d computed", res2[2])
	}
	if res2[3] != 5 {
		t.Errorf("error, result for 3 should be 5, %d computed", res2[3])
	}
}

func TestMajoritySCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := MajoritySCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestApprovalSWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 3, 2},
		{2, 3, 1},
	}
	thresholds := []int{2, 1, 2}

	res, _ := ApprovalSWF(prefs, thresholds)

	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 2 {
		t.Errorf("error, result for 2 should be 2, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestApprovalSCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
	}
	thresholds := []int{2, 1, 2}

	res, err := ApprovalSCF(prefs, thresholds)

	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestCondorcetWinner(t *testing.T) {
	prefs1 := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	prefs2 := [][]Alternative{
		{1, 2, 3},
		{2, 3, 1},
		{3, 1, 2},
	}

	prefs3 := [][]Alternative{
		{1, 2, 3},
		{1, 3, 2},
		{3, 2, 1},
	}

	prefs4 := [][]Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 3, 2, 4},
		{1, 3, 2, 4},
		{1, 3, 2, 4},
		{1, 3, 2, 4},
		{4, 2, 1, 3},
		{4, 2, 1, 3},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{4, 2, 3, 1},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{3, 2, 1, 4},
		{4, 3, 2, 1},
		{4, 3, 2, 1},
	}

	res1, _ := CondorcetWinner(prefs1)
	res2, _ := CondorcetWinner(prefs2)
	res3, _ := CondorcetWinner(prefs3)
	res4, _ := CondorcetWinner(prefs4)

	if len(res1) == 0 || res1[0] != 1 {
		t.Errorf("error, 1 should be the only best alternative for prefs1")
	}
	if len(res2) != 0 {
		t.Errorf("no best alternative for prefs2")
	}
	if len(res3) == 0 || res3[0] != 1 {
		t.Errorf("error, 1 should be the only best alternative for prefs3")
	}
	if len(res4) == 0 || res4[0] != 3 {
		t.Errorf("error, 3 should be the only best alternative for prefs4")
	}
}

func TestTieBreak(t *testing.T) {
	p := Profile{
		{2, 1, 3, 4, 5, 6},
		{5, 4, 2, 3, 1, 6},
		{5, 2, 3, 4, 1, 6},
		{2, 1, 3, 4, 5, 6},
		{2, 4, 3, 5, 1, 6},
		{5, 2, 3, 6, 4, 1},
	}

	orderedAlts := []Alternative{5, 2, 1, 3, 4, 6}

	fmt.Println("\nCount sans tie-break :")
	countWithoutTieBreak, _ := MajoritySWF(p)
	fmt.Println(countWithoutTieBreak)

	/*fmt.Println("\nRésultat avec tie-break simple:")
	fmt.Println(SWFFactory(MajoritySWF, TieBreak)(p))*/

	fmt.Println("\nRésultat avec tie-break factory :")
	fmt.Println(SWFFactory(MajoritySWF, TieBreakFactory(orderedAlts))(p))
}

func TestCopelandWinner(t *testing.T) {
	prefs1 := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	prefs2 := [][]Alternative{
		{1, 2, 3},
		{2, 3, 1},
		{3, 1, 2},
	}

	prefs3 := [][]Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
	}

	res1, _ := CopelandWinner(prefs1)
	res2, _ := CopelandWinner(prefs2)
	res3, _ := CopelandWinner(prefs3)

	if len(res1) == 0 || res1[0] != 1 {
		t.Errorf("error, 1 should be the only best alternative for prefs1")
	}
	if len(res2) != 3 {
		t.Errorf("3 bests alternative for prefs2")
	}
	if len(res3) != 2 {
		t.Errorf("error, 2 best alternatives for prefs3")
	}
}
