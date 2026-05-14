package diff

import "sort"

// Mismatch holds a key whose value differs between two env maps.
type Mismatch struct {
	Key      string
	LeftVal  string
	RightVal string
}

// Result holds the full comparison output between two env maps.
type Result struct {
	MissingInRight []string
	MissingInLeft  []string
	Mismatched     []Mismatch
}

// IsClean returns true when there are no differences.
func (r Result) IsClean() bool {
	return len(r.MissingInRight) == 0 &&
		len(r.MissingInLeft) == 0 &&
		len(r.Mismatched) == 0
}

// Compare compares two env maps and returns a Result describing differences.
func Compare(left, right map[string]string) Result {
	var result Result

	for k, lv := range left {
		if rv, ok := right[k]; !ok {
			result.MissingInRight = append(result.MissingInRight, k)
		} else if lv != rv {
			result.Mismatched = append(result.Mismatched, Mismatch{
				Key:      k,
				LeftVal:  lv,
				RightVal: rv,
			})
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, k)
		}
	}

	sortStrings(result.MissingInRight)
	sortStrings(result.MissingInLeft)
	sortMismatched(result.Mismatched)

	return result
}

func sortStrings(s []string) {
	sort.Strings(s)
}

func sortMismatched(m []Mismatch) {
	sort.Slice(m, func(i, j int) bool {
		return m[i].Key < m[j].Key
	})
}
