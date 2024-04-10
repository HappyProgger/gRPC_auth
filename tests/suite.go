package main

import "github.com/stretchr/testify/suite"

type Suite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}
