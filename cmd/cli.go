package main

import (
	"flag"
	"github.com/hazcod/aikido-sdk-go/config"
	"github.com/hazcod/aikido-sdk-go/pkg/aikido"
	"github.com/sirupsen/logrus"
)

func main() {
	//ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "config.yml", "The YAML configuration file.")
	flag.Parse()

	conf := config.Config{}
	if err := conf.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := conf.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.SetLevel(logrusLevel)

	// ---

	aikClient, err := aikido.New(logger, conf.Aikido.ClientID, conf.Aikido.ClientSecret)
	if err != nil {
		logger.WithError(err).Fatal("failed to create aikido client")
	}

	issues, err := aikClient.GetIssues(true)
	if err != nil {
		logger.WithError(err).Fatal("failed to get issues")
	}

	logger.Infof("Retrieved %d issues:", len(issues))
	openIssues := 0

	for _, issue := range issues {
		openIssues += 1
		logger.Info(issue.GetName())
	}

	logger.Infof("Retrieved %d open issues.", openIssues)
}
