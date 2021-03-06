package ecs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type TaskDefinition struct {
	Overwrite       bool
	DeleteContainer bool

	Family string

	TaskRoleArn      *string
	ExecutionRoleArn *string

	NetworkMode             *string
	ContainerDefinitions    []*ContainerDefinition
	Volumes                 []*ecs.Volume
	RequiresCompatibilities []*string

	Cpu    *string
	Memory *string

	IpcMode *string
	PidMode *string

	PlacementConstraints []*ecs.TaskDefinitionPlacementConstraint
	ProxyConfiguration   *ecs.ProxyConfiguration

	Tags []*ecs.Tag
}

func (td *TaskDefinition) isValid() error {
	if td.Family == "" {
		return fmt.Errorf("Task Definition must have a name")
	}
	if _, err := parseFamily(td.Family); err != nil {
		return fmt.Errorf("Task Definition Family cannot be parsed")
	}
	if td.ContainerDefinitions != nil {
		if len(td.ContainerDefinitions) < 1 {
			return fmt.Errorf("Container Definitions must have at least 1 Container")
		}
		for _, cd := range td.ContainerDefinitions {
			if cd == nil {
				return fmt.Errorf("Container Definitions cannot have nil value")
			} else {
				if err := cd.isValid(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (td *TaskDefinition) Register(svc *ecs.ECS) (*ecs.TaskDefinition, error) {
	if err := td.isValid(); err != nil {
		return nil, err
	}

	var taskDefinition *ecs.TaskDefinition
	if !td.Overwrite {
		tdout, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{TaskDefinition: &td.Family})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() != ecs.ErrCodeClientException {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		taskDefinition = tdout.TaskDefinition
	}

	if taskDefinition != nil && td.isEmpty() {
		fmt.Println("No changes were found, using the latest version of the Task Definition")
		return taskDefinition, nil
	}

	fmt.Printf("Registering new Task Definition from [%s]...\n", td.Family)
	input := td.generateInput(taskDefinition)
	tdnew, err := svc.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	version, _ := parseFamilyRevision(*tdnew.TaskDefinition.TaskDefinitionArn)
	fmt.Printf("Successfully registered [%s]\n\n", version)
	return tdnew.TaskDefinition, nil
}

func (td *TaskDefinition) Update(svc *ecs.ECS) (*ecs.TaskDefinition, error) {
	if err := td.isValid(); err != nil {
		return nil, err
	}

	if td.Overwrite {
		return nil, fmt.Errorf("Update cannot be run with Overwrite option, since this can be dangerous")
	}

	tdout, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{TaskDefinition: &td.Family})
	if err != nil {
		return nil, err
	}

	if tdout.TaskDefinition == nil {
		return nil, fmt.Errorf("Existing Task Definition was not found, cannot update")
	}

	if td.isEmpty() {
		fmt.Println("No changes were found, using the latest version of the Task Definition")
		return tdout.TaskDefinition, nil
	}

	fmt.Printf("Registering new Task Definition from [%s]...\n", td.Family)
	input := td.generateInput(tdout.TaskDefinition)
	tdnew, err := svc.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	version, _ := parseFamilyRevision(*tdnew.TaskDefinition.TaskDefinitionArn)
	fmt.Printf("Successfully registered [%s]\n\n", version)
	return tdnew.TaskDefinition, nil
}

func (td *TaskDefinition) isEmpty() bool {
	return !td.Overwrite &&
		td.TaskRoleArn == nil &&
		td.ExecutionRoleArn == nil &&
		td.NetworkMode == nil &&
		td.ContainerDefinitions == nil &&
		td.Volumes == nil &&
		td.RequiresCompatibilities == nil &&
		td.Cpu == nil &&
		td.Memory == nil &&
		td.IpcMode == nil &&
		td.PidMode == nil &&
		td.PlacementConstraints == nil &&
		td.ProxyConfiguration == nil &&
		td.Tags == nil
}

func (td *TaskDefinition) generateInput(old *ecs.TaskDefinition) *ecs.RegisterTaskDefinitionInput {
	taskInput := &ecs.RegisterTaskDefinitionInput{}

	family, _ := parseFamily(td.Family)
	taskInput.Family = &family

	if td.TaskRoleArn != nil {
		taskInput.TaskRoleArn = td.TaskRoleArn
	} else {
		taskInput.TaskRoleArn = old.TaskRoleArn
	}
	if td.ExecutionRoleArn != nil {
		taskInput.ExecutionRoleArn = td.ExecutionRoleArn
	} else {
		taskInput.ExecutionRoleArn = old.ExecutionRoleArn
	}
	if td.NetworkMode != nil {
		taskInput.NetworkMode = td.NetworkMode
	} else {
		taskInput.NetworkMode = old.NetworkMode
	}
	if td.ContainerDefinitions != nil {
		containerDefinitions := []*ecs.ContainerDefinition{}
		for i, cd := range td.ContainerDefinitions {
			containerDefinitions = append(containerDefinitions, cd.generateDefinition(old.ContainerDefinitions[i]))
		}
		taskInput.ContainerDefinitions = containerDefinitions
	} else {
		taskInput.ContainerDefinitions = old.ContainerDefinitions
	}
	if td.Volumes != nil {
		taskInput.Volumes = td.Volumes
	} else {
		taskInput.Volumes = old.Volumes
	}
	if td.RequiresCompatibilities != nil {
		taskInput.RequiresCompatibilities = td.RequiresCompatibilities
	} else {
		taskInput.RequiresCompatibilities = old.RequiresCompatibilities
	}
	if td.Cpu != nil {
		taskInput.Cpu = td.Cpu
	} else {
		taskInput.Cpu = old.Cpu
	}
	if td.Memory != nil {
		taskInput.Memory = td.Memory
	} else {
		taskInput.Memory = old.Memory
	}
	if td.IpcMode != nil {
		taskInput.IpcMode = td.IpcMode
	} else {
		taskInput.IpcMode = old.IpcMode
	}
	if td.PidMode != nil {
		taskInput.PidMode = td.PidMode
	} else {
		taskInput.PidMode = old.PidMode
	}
	if td.PlacementConstraints != nil {
		taskInput.PlacementConstraints = td.PlacementConstraints
	} else {
		taskInput.PlacementConstraints = old.PlacementConstraints
	}
	if td.ProxyConfiguration != nil {
		taskInput.ProxyConfiguration = td.ProxyConfiguration
	} else {
		taskInput.ProxyConfiguration = old.ProxyConfiguration
	}

	return taskInput
}

var arnRegex, _ = regexp.Compile(`^arn:aws:ecs:[a-z]{2}-[a-z]+-\d{1,2}:\d{12}:task-definition\/[\w-]+:\d+$`)
var familyRegex, _ = regexp.Compile(`^[\w-]+$`)
var familyRevisionRegex, _ = regexp.Compile(`^[\w-]+:\d+$`)

func parseFamily(family string) (string, error) {
	switch true {
	case arnRegex.MatchString(family):
		return strings.Split(strings.Split(family, ":")[5], "/")[1], nil
	case familyRegex.MatchString(family):
		return family, nil
	case familyRevisionRegex.MatchString(family):
		return strings.Split(family, ":")[0], nil
	}

	return "", fmt.Errorf("Family not found")
}

func parseFamilyRevision(family string) (string, error) {
	switch true {
	case arnRegex.MatchString(family):
		return strings.SplitN(family, "/", 2)[1], nil
	case familyRegex.MatchString(family):
		return family, nil
	case familyRevisionRegex.MatchString(family):
		return family, nil
	}

	return "", fmt.Errorf("Family not found")
}
