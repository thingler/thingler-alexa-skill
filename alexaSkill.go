package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	alexa "github.com/ericdaugherty/alexa-skills-kit-golang"
)

// Thingler handles requests from the Alexa skill.
type Thingler struct {
	config *Config
}

// Handle processes calls from Lambda
func Handle(ctx context.Context, requestEnv *alexa.RequestEnvelope) (interface{}, error) {

	// Get configurations from file
	config := &Config{}
	err := config.Parse()
	if err != nil {
		return nil, err
	}

	a := &alexa.Alexa{
		ApplicationID: config.ApplicationID,
		RequestHandler: &Thingler{
			config: config,
		},
		IgnoreApplicationID: false,
		IgnoreTimestamp:     false,
	}
	return a.ProcessRequest(ctx, requestEnv)
}

// OnSessionStarted called when a new session is created.
func (thing *Thingler) OnSessionStarted(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	log.Printf("OnSessionStarted requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	return nil
}

// OnLaunch called when a request is received of type LaunchRequest
func (thing *Thingler) OnLaunch(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {
	speechText := "Welcome to Thingler, you can ask me to turn the Thingler smart plug on or to turn the Thingler smart plug off"

	log.Printf("OnLaunch requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	response.SetSimpleCard(thing.config.CardTitle, speechText)
	response.SetOutputText(speechText)
	response.SetRepromptText(speechText)

	response.ShouldSessionEnd = false

	return nil
}

// OnIntent called when a request is received of type IntentRequest
func (thing *Thingler) OnIntent(ctx context.Context, request *alexa.Request, alexaSession *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	log.Printf("OnIntent requestId=%s, sessionId=%s, intent=%s", request.RequestID, alexaSession.SessionID, request.Intent.Name)

	sess := session.Must(session.NewSession())

	// Try first with Environment variables and secondly with IAM role
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})

	config := &aws.Config{
		Region:      &thing.config.Region,
		Credentials: creds,
		Endpoint:    &thing.config.IOTEndpoint,
	}

	clientIOT := iotdataplane.New(sess, config)

	turnOn := &IntentTurnOn{
		IOTClient: clientIOT,
		Response:  response,
		CardTitle: &thing.config.CardTitle,
		Topic:     &thing.config.IOTTopic,
	}

	turnOff := &IntentTurnOff{
		IOTClient: clientIOT,
		Response:  response,
		CardTitle: &thing.config.CardTitle,
		Topic:     &thing.config.IOTTopic,
	}

	help := &IntentHelp{
		Response:  response,
		CardTitle: &thing.config.CardTitle,
	}

	intent, err := NewIntentFactory().
		AddIntent(turnOn).
		AddIntent(turnOff).
		AddIntent(help).
		GetIntent(&request.Intent.Name)
	if err != nil {
		return err
	}

	err = intent.Do()

	return err
}

// OnSessionEnded called when a request is received of type SessionEndedRequest
func (thing *Thingler) OnSessionEnded(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	log.Printf("OnSessionEnded requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	return nil
}

func main() {
	lambda.Start(Handle)
}
