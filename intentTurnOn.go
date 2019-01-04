package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/iotdataplane"
	alexa "github.com/ericdaugherty/alexa-skills-kit-golang"
)

// IntentTurnOn
type IntentTurnOn struct {
	IOTClient *iotdataplane.IoTDataPlane
	Response  *alexa.Response
	CardTitle *string
	Topic     *string
}

// Name return the intent name
func (intent *IntentTurnOn) Name() string {
	return "TurnOnIntent"
}

// Do perform the Turn On intent
func (intent *IntentTurnOn) Do() error {

	log.Printf("%s triggered", intent.Name())

	speechText := "Turning on Thingler Smart Plug"

	publishInput := &iotdataplane.PublishInput{
		Topic:   intent.Topic,
		Payload: []byte("on"),
	}

	_, err := intent.IOTClient.Publish(publishInput)
	if err != nil {
		return err
	}

	intent.Response.SetSimpleCard(*intent.CardTitle, speechText)
	intent.Response.SetOutputText(speechText)

	log.Printf("Set Output speech, value now: %s", intent.Response.OutputSpeech.Text)
	return err
}
