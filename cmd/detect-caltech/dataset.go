package main

import "fmt"

type Dataset struct {
	Dir  string
	Skip int
	Ext  string
	Seqs map[int][]int
}

func datasetByName(name string) (*Dataset, error) {
	switch name {
	case "usa":
		return &Dataset{"USA", 30, "jpg", map[int][]int{
			0:  interval(0, 15),
			1:  interval(0, 6),
			2:  interval(0, 12),
			3:  interval(0, 13),
			4:  interval(0, 12),
			5:  interval(0, 13),
			6:  interval(0, 19),
			7:  interval(0, 12),
			8:  interval(0, 11),
			9:  interval(0, 12),
			10: interval(0, 12),
		}}, nil
	case "usatrain":
		return &Dataset{"USA", 30, "jpg", map[int][]int{
			0: interval(0, 15),
			1: interval(0, 6),
			2: interval(0, 12),
			3: interval(0, 13),
			4: interval(0, 12),
			5: interval(0, 13),
		}}, nil
	case "usatest":
		return &Dataset{"USA", 30, "jpg", map[int][]int{
			6:  interval(0, 19),
			7:  interval(0, 12),
			8:  interval(0, 11),
			9:  interval(0, 12),
			10: interval(0, 12),
		}}, nil
	default:
		return nil, fmt.Errorf(`unknown dataset: "%s"`, name)
	}
}

func interval(a, b int) []int {
	x := make([]int, b-a)
	for i := a; i < b; i++ {
		x[i-a] = i
	}
	return x
}
