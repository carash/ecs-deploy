package ecs

import (
	"fmt"
	"time"

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

func (p *ServicePlugin) UpdateService(timeout int64) error {
	svc := ecs.New(p.AWSCredential.newSession())
	service, err := p.Service.Update(svc)
	if err != nil {
		return err
	}

	start := time.Now()
	check := make(chan error)
	td, _ := parseFamilyRevision(*service.TaskDefinition)

	go func() {
		for {
			go func() {
				taskout, err := svc.ListTasks(&ecs.ListTasksInput{
					Cluster:     p.Service.Cluster,
					ServiceName: &p.Service.Service,
				})
				if err != nil {
					check <- err
				}
				if len(taskout.TaskArns) == 0 {
					return
				}

				detout, err := svc.DescribeTasks(&ecs.DescribeTasksInput{
					Cluster: p.Service.Cluster,
					Tasks:   taskout.TaskArns,
				})
				if err != nil {
					check <- err
				}

				healthy := int64(0)
				for _, t := range detout.Tasks {
					taskDefinition, _ := parseFamilyRevision(*t.TaskDefinitionArn)
					fmt.Printf("Status of [%s] -> %s\n", taskDefinition, *t.HealthStatus)

					if *t.TaskDefinitionArn == *service.TaskDefinition && *t.HealthStatus == "HEALTHY" {
						healthy += 1
					}
				}
				fmt.Println()

				if int64(len(taskout.TaskArns)) != *service.DesiredCount {
					return
				}
				if healthy == *service.DesiredCount {
					fmt.Printf("Task [%s] is HEALTHY after %d seconds\n\n", td, int64(time.Now().Sub(start).Seconds()))
					check <- nil
					return
				}
			}()

			time.Sleep(10 * time.Second)
			fmt.Printf("Waiting for Task [%s] to be HEALTHY, %ds...\n", td, int64(time.Now().Sub(start).Seconds()))
		}
	}()

	select {
	case err := <-check:
		if err != nil {
			return err
		}
	case <-time.After(time.Duration(timeout) * time.Second):
		return fmt.Errorf("Timed out after %ds while wating for Task [%s] to deploy\n\n", int64(time.Now().Sub(start).Seconds()), td)
	}

	return nil

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
