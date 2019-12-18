package cmd

import (
	"fmt"
	"testing"
)

func Test_createNewTarget(t *testing.T) {
	result, err := createNewTarget("bloodAndWine", []string{"good", "call"})
	fmt.Println(result)
	fmt.Println(err)
}
