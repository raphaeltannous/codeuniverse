package main

import (
	"reflect"
	"testing"
)

func equalSlices(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	visited := make([]bool, len(b))
	for _, arrA := range a {
		found := false
		for j, arrB := range b {
			if !visited[j] && reflect.DeepEqual(arrA, arrB) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestPermute(t *testing.T) {
	cases := []struct {
		nums []int
		want [][]int
	}{
		{[]int{1, 2, 3}, [][]int{
			{}, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3},
		}},
		{[]int{0}, [][]int{
			{}, {0},
		}},
	}

	for _, c := range cases {
		got := subsets(c.nums)
		if !equalSlices(got, c.want) {
			t.Errorf("subsets(%v) = %v, want %v", c.nums, got, c.want)
		}
	}
}
