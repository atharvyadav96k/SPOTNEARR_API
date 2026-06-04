package main

import (
	"fmt"

	"github.com/atharvyadav96k/SPOTNEARR_API/factories/token"
)

type searchToken struct {
	brandNameToken []string
	descToken      []string
}

func GenerateTokens(val string) *searchToken {
	result := token.QueryTokenParse(val)
	fmt.Println(result)
	return &searchToken{}
}
