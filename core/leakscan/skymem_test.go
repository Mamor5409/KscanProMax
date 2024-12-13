package leakscan

import (
	"fmt"
	"testing"
)

func TestSearchSkymem(t *testing.T) {
	strcode := "baidu.com"
	a := SearchSkymem(strcode)
	fmt.Println(a)
}
