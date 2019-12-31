package filereader

import (
	"fmt"
	"testing"
)

func TestParsePercentage(t *testing.T) {
	var v float64
	var err error

	v, err = parsePercentage("1.5")
	if err != nil || v != 1.5 {
		fmt.Println(v, err, v != 1.5)
		t.Error("happy path failed")
	}

	v, err = parsePercentage("0")
	if err != nil || v != 0 {
		fmt.Println(v, err, v != 0)
		t.Error("happy path (0) failed")
	}

	v, err = parsePercentage("100.00")
	if err != nil || v != 100 {
		fmt.Println(v, err, v != 100)
		t.Error("happy path (100) failed")
	}

	v, err = parsePercentage("-0.5")
	if err.Error() != "incorrect percentage value" {
		t.Error("negative value failed")
	}

	v, err = parsePercentage("100.01")
	if err.Error() != "incorrect percentage value" {
		t.Error("too large value failed")
	}
}
