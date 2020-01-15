package ksyun

import (
	"fmt"
	"testing"
)

func TestHump2Downline(t *testing.T) {
	ss := Hump2Downline("HsfjdSjkfjdkD")
	fmt.Println(ss)
}
func TestDownline2Hump(t *testing.T) {
	ss := Downline2Hump("hsfjd_sjkfjdk_d")
	fmt.Println(ss)
}
