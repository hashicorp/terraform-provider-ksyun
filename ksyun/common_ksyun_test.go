package ksyun

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
func TestT(t *testing.T) {
	a := assert.New(t)
	a.Equal(Camel2Hungarian("SupportIPV6"), "support_ipv6")
	a.Equal(Camel2Hungarian("DBInstanceIdentifier"), "db_instance_identifier")
	a.Equal(Camel2Hungarian("AA.1.DBInstanceIdentifier"), "aa_1.db_instance_identifier")
}
