package credential

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Credential struct {
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSAssumeRoleARN   string
	AWSRegion          string
}

func (c *Credential) NewSession() *session.Session {
	awsConfig := aws.Config{}

	if c.AWSAccessKeyID != "" && c.AWSSecretAccessKey != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(c.AWSAccessKeyID, c.AWSSecretAccessKey, "")
	} else if c.AWSAssumeRoleARN != "" {
		awsConfig.Credentials = stscreds.NewCredentials(session.Must(session.NewSession()), c.AWSAssumeRoleARN)
	}
	if c.AWSRegion != "" {
		awsConfig.Region = aws.String(c.AWSRegion)
	}

	return session.Must(session.NewSession(&awsConfig))
}
