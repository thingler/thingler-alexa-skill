package main

import (
	"log"
	"regexp"

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

	speechText := `<speak>You can ask the <phoneme alphabet='ipa' ph='Thiŋler'>Thingler</phoneme>
smart plug to turn on, or to turn off</speak>`

	cleanText := regexp.MustCompile("<[^>]*>").
		ReplaceAllString(speechText, "")

	intent.Response.SetSimpleCard(*intent.CardTitle, cleanText)
	intent.Response.SetOutputSSML(speechText)

	return nil
}
