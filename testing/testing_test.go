package testing

import "testing"

func TestRandStr(t *testing.T) {
	for i := 1; i <= 30; i++ {
		str1 := RandStr(i)
		str2 := RandStr(i)

		if str1 == str2 {
			t.Errorf("%s should not equal %s", str1, str2)
		}
	}
}
