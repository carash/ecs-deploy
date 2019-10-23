package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Credential struct {
	AWSAccessKeyID     *string
	AWSSecretAccessKey *string
	AWSAssumeRoleARN   *string
	AWSRegion          *string
}

type ServicePlugin struct {
	AWSCredential Credential
	Service       Service
}

type TaskPlugin struct {
	AWSCredential  Credential
	TaskDefinition TaskDefinition
}

func (c *Credential) newSession() *session.Session {
	awsConfig := aws.Config{}

	if c.AWSAccessKeyID != nil && c.AWSSecretAccessKey != nil {
		awsConfig.Credentials = credentials.NewStaticCredentials(*c.AWSAccessKeyID, *c.AWSSecretAccessKey, "")
	} else if c.AWSAssumeRoleARN != nil {
		awsConfig.Credentials = stscreds.NewCredentials(session.Must(session.NewSession()), *c.AWSAssumeRoleARN)
	}
	if c.AWSRegion != nil {
		awsConfig.Region = aws.String(*c.AWSRegion)
	}

	return session.Must(session.NewSession(&awsConfig))
}

func (p *ServicePlugin) DeployService() error {
	svc := ecs.New(p.AWSCredential.newSession())
	_, err := p.Service.Deploy(svc)
	return err
}

func (p *ServicePlugin) UpdateService() error {
	svc := ecs.New(p.AWSCredential.newSession())
	_, err := p.Service.Update(svc)
	return err
}

func (p *TaskPlugin) RegisterTask() error {
	svc := ecs.New(p.AWSCredential.newSession())
	_, err := p.TaskDefinition.Register(svc)
	return err
}

func (p *TaskPlugin) UpdateTask() error {
	svc := ecs.New(p.AWSCredential.newSession())
	_, err := p.TaskDefinition.Update(svc)
	return err
}
