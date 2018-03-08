package testing

import "testing"

func TestRandStr(t *testing.T) {
	length := 16
	for i := 0; i <= 40; i++ {
		str1 := RandStr(length)
		str2 := RandStr(length)

		if len(str1) != length {
			t.Errorf("%s should have %d characters", str1, length)
		}
		if len(str2) != length {
			t.Errorf("%s should have %d characters", str2, length)
		}
		if str1 == str2 {
			t.Errorf("%s should not equal %s", str1, str2)
		}
	}
}
