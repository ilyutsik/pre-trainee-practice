package pipeline

import (
	"slices"
	"testing"
)

func TestDouble(t *testing.T) {
	cases := map[string]struct {
		in   []int
		want []int
	}{
		"empty in slice":             {[]int{}, []int{}},
		"nil in slice":               {nil, nil},
		"zero in slice":              {[]int{0}, []int{0}},
		"simple in slice":            {[]int{1, 2, 3, 4, 5}, []int{2, 4, 6, 8, 10}},
		"in slice with negative int": {[]int{-11, 0, 39, -45, -50}, []int{-22, 0, 78, -90, -100}}}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := Double(tc.in)
			if (got == nil) != (tc.want == nil) {
				t.Errorf("%v: get %#v want %#v", name, got, tc.want)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("%q, got %v want %v", name, got, tc.want)
			}
		})
	}
}

func TestDoubleAliasing(t *testing.T) {
	base := []int{1, 2, 3, 4, 5}
	orig := slices.Clone(base)
	got := Double(base[:2])
	got[0] = 99
	if !slices.Equal(base, orig) {
		t.Errorf("Double func create aliasing on input and out vale: base %v orig %v", base, orig)
	}
}

func TestFilterMin(t *testing.T) {
	cases := map[string]struct {
		in   []int
		want []int
		min  int
	}{
		"nil in slice":                         {nil, nil, 12},
		"empty in slice":                       {[]int{}, []int{}, 0},
		"all slice values < min":               {[]int{1, -1, 2}, []int{}, 3},
		"in slice with zero min value":         {[]int{1, -5, 15, 7, -9, 8, -65}, []int{1, 15, 7, 8}, 0},
		"in slice with positive(5) min value":  {[]int{1, -5, 15, 7, -9, 8, -65}, []int{15, 7, 8}, 5},
		"in slice with negative(-5) min value": {[]int{1, -5, 15, 7, -9, 8, -65}, []int{1, -5, 15, 7, 8}, -5},
		"in slice with el = min value":         {[]int{3, 3, 3}, []int{3, 3, 3}, 3}}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := FilterMin(c.in, c.min)
			if (got == nil) != (c.want == nil) {
				t.Errorf("%v: get %#v want %#v", name, got, c.want)
			}
			if !slices.Equal(got, c.want) {
				t.Errorf("%v: get %v want %v", name, got, c.want)
			}
		})
	}
}

func TestFilterMinAliasing(t *testing.T) {
	base := []int{5, 6, 1, 2, 3, 4, 5}
	orig := slices.Clone(base)
	got := FilterMin(base[:2], 2)
	got[0] = 99
	if !slices.Equal(base, orig) {
		t.Errorf("FilterMin func create aliasing on input and out vale: base %v orig %v", base, orig)
	}
}

func TestUniqueStable(t *testing.T) {
	cases := map[string]struct {
		in   []int
		want []int
	}{
		"nil in slice":                      {nil, nil},
		"empty in slice":                    {[]int{}, []int{}},
		"in slice with one value":           {[]int{12}, []int{12}},
		"in slice with duplicates":          {[]int{1, 2, 3, 2, 4, 2, 4, 1, 5, 2}, []int{1, 2, 3, 4, 5}},
		"in slice with negative duplicates": {[]int{1, -1, 2, 1, -1, -2, 2, -2}, []int{1, -1, 2, -2}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := UniqueStable(c.in)
			if (got == nil) != (c.want == nil) {
				t.Errorf("%v: get %#v want %#v", name, got, c.want)
			}
			if !slices.Equal(got, c.want) {
				t.Errorf("%v: got %v want %v", name, got, c.want)
			}
		})
	}
}

func TestUniqueStablesAliasing(t *testing.T) {
	base := []int{1, 2, 2, 3, 4, 5, 5}
	orig := slices.Clone(base)
	got := UniqueStable(base[:5])
	got[0] = 99
	if !slices.Equal(base, orig) {
		t.Errorf("UniqueStable func create aliasing on input and out vale: base %v orig %v", base, orig)
	}
}
