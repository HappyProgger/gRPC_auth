package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Hello(name string, language string) (string, error) {
	if name == "" {
		name = "World"
	}

	prefix := ""

	switch language {
	case "english":
		prefix = "Hello"
	case "spanish":
		prefix = "Hola"
	case "german":
		prefix = "Hallo"
	default:
		return "", errors.New("need to provide a supported language")
	}

	return prefix + " " + name, nil
}

func TestEnglish(t *testing.T) {
	name := "Ben"
	language := "english"
	expected := "Hello Ben"
	actual, err := Hello(name, language)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
