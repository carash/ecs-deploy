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
	ContainerDefinitions    *[]*ContainerDefinition
	Volumes                 *[]*ecs.Volume
	RequiresCompatibilities *[]*string

	Cpu    *string
	Memory *string

	IpcMode *string
	PidMode *string

	PlacementConstraints *[]*ecs.TaskDefinitionPlacementConstraint
	ProxyConfiguration   *ecs.ProxyConfiguration

	Tags *[]*ecs.Tag
}

var defaultTaskMemory = "1024"

func (td *TaskDefinition) isValid() error {
	if td.Family == "" {
		return fmt.Errorf("Task Definition must have a name")
	}
	if _, err := parseFamily(td.Family); err != nil {
		return fmt.Errorf("Task Definition Family cannot be parsed")
	}
	if td.ContainerDefinitions != nil {
		if len(*td.ContainerDefinitions) < 1 {
			return fmt.Errorf("Container Definitions must have at least 1 Container")
		}
		for _, cd := range *td.ContainerDefinitions {
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
	input := td.parse(taskDefinition).unpackRegisterInput()
	tdnew, err := svc.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Successfully registered [%s]\n", *tdnew.TaskDefinition.TaskDefinitionArn)

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
	input := td.parse(tdout.TaskDefinition).unpackRegisterInput()
	tdnew, err := svc.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Successfully registered [%s]\n", *tdnew.TaskDefinition.TaskDefinitionArn)

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

func (td *TaskDefinition) parse(taskDefinition *ecs.TaskDefinition) *TaskDefinition {
	if taskDefinition == nil {
		return td
	}

	if td.TaskRoleArn == nil {
		td.TaskRoleArn = taskDefinition.TaskRoleArn
	}
	if td.ExecutionRoleArn == nil {
		td.ExecutionRoleArn = taskDefinition.ExecutionRoleArn
	}
	if td.NetworkMode == nil {
		td.NetworkMode = taskDefinition.NetworkMode
	}
	if td.ContainerDefinitions == nil {
		containers := []*ContainerDefinition{}
		for _, cd := range taskDefinition.ContainerDefinitions {
			containers = append(containers, (&ContainerDefinition{Name: *cd.Name}).parse(cd))
		}
		td.ContainerDefinitions = &containers
	} else {
		// map avaiable containers
		containerMap := make(map[string]*ecs.ContainerDefinition)
		for _, cd := range taskDefinition.ContainerDefinitions {
			containerMap[*cd.Name] = cd
		}

		// parse changes for provided containers
		containers := []*ContainerDefinition{}
		for _, cd := range *td.ContainerDefinitions {
			containers = append(containers, cd.parse(containerMap[cd.Name]))
			delete(containerMap, cd.Name)
		}

		// ignore non explicit containers
		if !td.DeleteContainer {
			for _, cd := range containerMap {
				containers = append(containers, (&ContainerDefinition{Name: *cd.Name}).parse(cd))
			}
		}

		td.ContainerDefinitions = &containers
	}
	if td.Volumes == nil {
		td.Volumes = &taskDefinition.Volumes
	}
	if td.RequiresCompatibilities == nil {
		td.RequiresCompatibilities = &taskDefinition.RequiresCompatibilities
	}
	if td.Cpu == nil {
		td.Cpu = taskDefinition.Cpu
	}
	if td.Memory == nil {
		td.Memory = taskDefinition.Memory
	}
	if td.IpcMode == nil {
		td.IpcMode = taskDefinition.IpcMode
	}
	if td.PidMode == nil {
		td.PidMode = taskDefinition.PidMode
	}
	if td.PlacementConstraints == nil {
		td.PlacementConstraints = &taskDefinition.PlacementConstraints
	}
	if td.ProxyConfiguration == nil {
		td.ProxyConfiguration = taskDefinition.ProxyConfiguration
	}

	return td
}

func (td *TaskDefinition) unpackRegisterInput() *ecs.RegisterTaskDefinitionInput {
	registerTaskDefinitionInput := &ecs.RegisterTaskDefinitionInput{}

	family, _ := parseFamily(td.Family)
	registerTaskDefinitionInput.Family = &family

	if td.TaskRoleArn != nil {
		registerTaskDefinitionInput.TaskRoleArn = td.TaskRoleArn
	}
	if td.ExecutionRoleArn != nil {
		registerTaskDefinitionInput.ExecutionRoleArn = td.ExecutionRoleArn
	}
	if td.NetworkMode != nil {
		registerTaskDefinitionInput.NetworkMode = td.NetworkMode
	}
	if td.ContainerDefinitions != nil {
		containerDefinitions := []*ecs.ContainerDefinition{}
		for _, cd := range *td.ContainerDefinitions {
			containerDefinitions = append(containerDefinitions, cd.unpack())
		}
		registerTaskDefinitionInput.ContainerDefinitions = containerDefinitions
	}
	if td.Volumes != nil {
		registerTaskDefinitionInput.Volumes = *td.Volumes
	}
	if td.RequiresCompatibilities != nil {
		registerTaskDefinitionInput.RequiresCompatibilities = *td.RequiresCompatibilities
	}
	if td.Cpu != nil {
		registerTaskDefinitionInput.Cpu = td.Cpu
	}
	if td.Memory != nil {
		registerTaskDefinitionInput.Memory = td.Memory
	} else {
		registerTaskDefinitionInput.Memory = &defaultTaskMemory
	}
	if td.IpcMode != nil {
		registerTaskDefinitionInput.IpcMode = td.IpcMode
	}
	if td.PidMode != nil {
		registerTaskDefinitionInput.PidMode = td.PidMode
	}
	if td.PlacementConstraints != nil {
		registerTaskDefinitionInput.PlacementConstraints = *td.PlacementConstraints
	}
	if td.ProxyConfiguration != nil {
		registerTaskDefinitionInput.ProxyConfiguration = td.ProxyConfiguration
	}

	return registerTaskDefinitionInput
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
