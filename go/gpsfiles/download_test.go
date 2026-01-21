package gpsfiles

import (
	"fmt"
	"testing"
)

func TestComputeFirstURL(t *testing.T) {
	fmt.Println(computeFirstURL(69, ACCURACY_1))
	fmt.Println(computeFirstURL(69, ACCURACY_5))
	fmt.Println(computeFirstURL(69, ACCURACY_25))
}

func TestFetchingDepartmentUrl(t *testing.T) {
	FetchingDepartmentUrl(69, ACCURACY_1)
	FetchingDepartmentUrl(69, ACCURACY_5)
	FetchingDepartmentUrl(69, ACCURACY_25)
}
