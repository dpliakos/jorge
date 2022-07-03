package jorge

import "testing"

func TestContains(t *testing.T) {
	if found := Contains([]string{"Hello", "there"}, "there"); found == false {
		t.Fatal()
	}

	if found := Contains([]string{"this", "is", "a", "sentence"}, "a"); found == false {
		t.Fatal()
	}

	if found := Contains([]string{"this", "is", "a", "sentence"}, "not"); found == true {
		t.Fatal()
	}
}
