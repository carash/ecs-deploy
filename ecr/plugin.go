package ecr

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type Credential struct {
	AWSAccessKeyID     *string
	AWSSecretAccessKey *string
	AWSAssumeRoleARN   *string
	AWSRegion          *string
}

type ImagePlugin struct {
	AWSCredential Credential
	Image         Image
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

func (p *ImagePlugin) FindImage() error {
	reg := ecr.New(p.AWSCredential.newSession())
	_, err := p.Image.Find(reg)
	return err
}

func (p *ImagePlugin) WaitForImage(interval, timeout int64) error {
	reg := ecr.New(p.AWSCredential.newSession())

	start := time.Now()
	delay := time.Duration(interval) * time.Second
	check := make(chan error)

	go func() {
		for {
			go func() {
				imgs, err := p.Image.Find(reg)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "ImageNotFoundException" {
						return
					}
					check <- err
					return
				}
				if len(imgs) > 0 {
					fmt.Printf("Image [%s] found after %d seconds\n", p.Image.DockerTag(), int64(time.Now().Sub(start).Seconds()))
					check <- nil
					return
				}
			}()

			time.Sleep(delay)
			fmt.Printf("Waiting for Image [%s], %ds...\n", p.Image.DockerTag(), int64(time.Now().Sub(start).Seconds()))
		}
	}()

	select {
	case err := <-check:
		if err != nil {
			return err
		}
	case <-time.After(time.Duration(timeout) * time.Second):
		return fmt.Errorf("Timed out after %ds while waiting for Image [%s]", int64(time.Now().Sub(start).Seconds()), p.Image.DockerTag())
	}

	return nil
}
