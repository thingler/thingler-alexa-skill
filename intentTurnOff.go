package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	alexa "github.com/ericdaugherty/alexa-skills-kit-golang"
)

// IntentTurnOff
type IntentTurnOff struct {
	IOTClient *iotdataplane.IoTDataPlane
	Response  *alexa.Response
	CardTitle *string
}

// Name return the intent name
func (intent *IntentTurnOff) Name() string {
	return "TurnOffIntent"
}

// Do perform the intent
func (intent *IntentTurnOff) Do() error {

	log.Printf("%s triggered", intent.Name())

	speechText := "Turning off Smart Plug"

	publishInput := &iotdataplane.PublishInput{
		Topic:   aws.String("tonttu/dimmer"),
		Payload: []byte("off"),
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
