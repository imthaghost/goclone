package auth

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
)

// Credentials gets username and password from user
func Credentials() (string, string) {
	qs := []*survey.Question{
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "What is your username?"},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Please enter your password."},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Username string
		Password string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatal(err)
	}

	return answers.Username, answers.Password
}
