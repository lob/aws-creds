package testing

import "testing"

func TestRandStr(t *testing.T) {
	for i := 0; i <= 40; i++ {
		str1 := RandStr(16)
		str2 := RandStr(16)

		if str1 == str2 {
			t.Errorf("%s should not equal %s", str1, str2)
		}
	}
}
