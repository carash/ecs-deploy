package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type Service struct {
	Cluster *string
	Service string

	PlatformVersion      *string
	NetworkConfiguration *ecs.NetworkConfiguration
	TaskDefinition       *TaskDefinition
	taskDefinition       *ecs.TaskDefinition

	ForceNewDeployment            *bool
	DeploymentConfiguration       *ecs.DeploymentConfiguration
	DesiredCount                  *int64
	HealthCheckGracePeriodSeconds *int64
}

func (s *Service) isValid() error {
	if s.Service == "" {
		return fmt.Errorf("Service must have a name")
	}

	return nil
}

func (s *Service) Deploy(svc *ecs.ECS) (*ecs.Service, error) {
	if err := s.isValid(); err != nil {
		return nil, err
	}

	// check availability of service
	srvout, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  s.Cluster,
		Services: []*string{&s.Service},
	})
	if err != nil {
		return nil, err
	}

	if len(srvout.Services) == 0 {
		return nil, fmt.Errorf("Cluster/Service combination not found")
	} else if len(srvout.Services) > 1 {
		return nil, fmt.Errorf("You can only update exactly 1 Service")
	}
	srv := srvout.Services[0]

	if s.TaskDefinition != nil {
		if s.TaskDefinition.Family == "" {
			s.TaskDefinition.Family = *srv.TaskDefinition
		}
		s.taskDefinition, err = s.TaskDefinition.Register(svc)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Deploying Service [%s]...\n", s.Service)
	input := s.unpackUpdateInput()
	snew, err := svc.UpdateService(input)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Successfully deployed [%s]\n\n", *snew.Service.ServiceName)
	return snew.Service, nil
}

func (s *Service) Update(svc *ecs.ECS) (*ecs.Service, error) {
	if err := s.isValid(); err != nil {
		return nil, err
	}

	// check availability of service
	srvout, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  s.Cluster,
		Services: []*string{&s.Service},
	})
	if err != nil {
		return nil, err
	}

	if len(srvout.Services) != 1 {
		return nil, fmt.Errorf("You can only update exactly 1 Service")
	}
	srv := srvout.Services[0]

	if s.TaskDefinition != nil {
		if s.TaskDefinition.Family == "" {
			s.TaskDefinition.Family = *srv.TaskDefinition
		} else {
			tdfam, err := parseFamily(s.TaskDefinition.Family)
			if err != nil {
				return nil, err
			}
			srvtdfam, err := parseFamily(*srv.TaskDefinition)
			if err != nil {
				return nil, err
			}
			if tdfam != srvtdfam {
				return nil, fmt.Errorf("You cannot change the task definition during and update operation")
			}
		}
		s.taskDefinition, err = s.TaskDefinition.Update(svc)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Updating Service [%s]...\n", s.Service)
	input := s.unpackUpdateInput()
	snew, err := svc.UpdateService(input)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Successfully updated [%s]\n\n", *snew.Service.ServiceName)
	return snew.Service, nil
}

func (s *Service) unpackUpdateInput() *ecs.UpdateServiceInput {
	updateServiceInput := &ecs.UpdateServiceInput{}

	updateServiceInput.Cluster = s.Cluster
	updateServiceInput.Service = &s.Service
	updateServiceInput.PlatformVersion = s.PlatformVersion
	updateServiceInput.NetworkConfiguration = s.NetworkConfiguration
	if s.taskDefinition != nil {
		updateServiceInput.TaskDefinition = s.taskDefinition.TaskDefinitionArn
	}
	updateServiceInput.ForceNewDeployment = s.ForceNewDeployment
	updateServiceInput.DeploymentConfiguration = s.DeploymentConfiguration
	updateServiceInput.DesiredCount = s.DesiredCount
	updateServiceInput.HealthCheckGracePeriodSeconds = s.HealthCheckGracePeriodSeconds

	return updateServiceInput
}
