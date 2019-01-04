package main

import (
	"log"

	alexa "github.com/ericdaugherty/alexa-skills-kit-golang"
)

// IntentHelp
type IntentHelp struct {
	Response  *alexa.Response
	CardTitle *string
}

// Name return the intent name
func (intent *IntentHelp) Name() string {
	return "AMAZON.HelpIntent"
}

// Do perform the intent
func (intent *IntentHelp) Do() error {

	log.Printf("%s triggered", intent.Name())

	speechText := "You can ask the Thingler smart plug to turn on or to turn off"

	intent.Response.SetSimpleCard(*intent.CardTitle, speechText)
	intent.Response.SetOutputText(speechText)
	intent.Response.SetRepromptText(speechText)

	return nil
}
