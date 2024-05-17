package aikido

import (
	"errors"
	"github.com/sirupsen/logrus"
)

type Aikido struct {
	logger    *logrus.Logger
	authToken *tokenGenerator
}

func New(logger *logrus.Logger, clientID, clientSecret string) (*Aikido, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	if clientID == "" || clientSecret == "" {
		return nil, errors.New("clientID or clientSecret required")
	}

	tokenGenerator, err := newTokenGenerator(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	return &Aikido{
		logger:    logger,
		authToken: tokenGenerator,
	}, nil
}
