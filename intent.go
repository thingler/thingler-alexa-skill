package main

import (
	"errors"
	"fmt"
)

type Intent interface {
	Name() string
	Do() error
}

// intentFactory contains a map[string] of all supported intents
// index value is the returned string from the intent's Name() function
type intentFactory struct {
	intents map[string]Intent
}

// NewIntentFactory is the constructor function for the intents
func NewIntentFactory() *intentFactory {
	factory := &intentFactory{}
	factory.intents = make(map[string]Intent)
	return factory
}

// Add intent
func (factory *intentFactory) AddIntent(intent Intent) *intentFactory {
	factory.intents[intent.Name()] = intent
	return factory
}

// Return the specified intent
func (factory *intentFactory) GetIntent(name *string) (Intent, error) {

	var err error

	intent, registered := factory.intents[*name]
	if !registered {
		errorMessage := fmt.Sprintf(
			"Intent, %s has not been implemented. Available intents:",
			*name,
		)
		for _, val := range factory.intents {
			errorMessage += fmt.Sprintf("\n %s", val.Name())
		}

		err = errors.New(errorMessage)
	}
	return intent, err
}
