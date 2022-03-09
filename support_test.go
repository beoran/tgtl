package tgtl

import "testing"
import "reflect"

func TestSort(t *testing.T) {
	arr := StringList("banana", "pear", "apple")
	expect := StringList("apple", "banana", "pear")
	sorted := arr.SortStrings()
	if !reflect.DeepEqual(sorted, expect) {
		t.Errorf("Not equal: %v<->%v", sorted, expect)
	}
	arr = StringList("banana")
	expect = StringList("banana")
	sorted = arr.SortStrings()
	if !reflect.DeepEqual(sorted, expect) {
		t.Errorf("Not equal: %v<->%v", sorted, expect)
	}
	arr = StringList()
	expect = StringList()
	sorted = arr.SortStrings()
	if !reflect.DeepEqual(sorted, expect) {
		t.Errorf("Not equal: %v<->%v", sorted, expect)
	}
}
