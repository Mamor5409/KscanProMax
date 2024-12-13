package leakscan

import (
	"fmt"
	"testing"
)

func TestSearchBySearchcodeAPI(t *testing.T) {
	strcode := "jnjkjt.com"
	a := SearchBySearchcodeApi(strcode)
	fmt.Println(a)
}
