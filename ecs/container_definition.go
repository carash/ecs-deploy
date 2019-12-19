package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type ContainerDefinition struct {
	Overwrite bool

	Name       string
	Image      *string
	EntryPoint []*string
	Command    []*string

	DependsOn             []*ecs.ContainerDependency
	RepositoryCredentials *ecs.RepositoryCredentials
	DockerSecurityOptions []*string
	DockerLabels          map[string]*string

	Essential            *bool
	Cpu                  *int64
	Memory               *int64
	MemoryReservation    *int64
	ResourceRequirements []*ecs.ResourceRequirement
	Ulimits              []*ecs.Ulimit

	User                   *string
	WorkingDirectory       *string
	Interactive            *bool
	PseudoTerminal         *bool
	ReadonlyRootFilesystem *bool
	Privileged             *bool
	LinuxParameters        *ecs.LinuxParameters
	SystemControls         []*ecs.SystemControl

	Hostname    *string
	ExtraHosts  []*ecs.HostEntry
	MountPoints []*ecs.MountPoint
	VolumesFrom []*ecs.VolumeFrom

	Environment []*ecs.KeyValuePair
	Secrets     []*ecs.Secret

	Links        []*string
	PortMappings []*ecs.PortMapping

	DisableNetworking *bool
	DnsSearchDomains  []*string
	DnsServers        []*string

	HealthCheck  *ecs.HealthCheck
	StartTimeout *int64
	StopTimeout  *int64

	FirelensConfiguration *ecs.FirelensConfiguration
	LogConfiguration      *ecs.LogConfiguration
}

func (cd *ContainerDefinition) isValid() error {
	if cd.Name == "" {
		return fmt.Errorf("Container Definitions must have a name")
	}

	return nil
}

func (cd *ContainerDefinition) generateDefinition(old *ecs.ContainerDefinition) *ecs.ContainerDefinition {
	container := &ecs.ContainerDefinition{}

	container.Name = &cd.Name

	if cd.Image != nil {
		container.Image = cd.Image
	} else {
		container.Image = old.Image
	}
	if cd.EntryPoint != nil {
		container.EntryPoint = cd.EntryPoint
	} else {
		container.EntryPoint = old.EntryPoint
	}
	if cd.Command != nil {
		container.Command = cd.Command
	} else {
		container.Command = old.Command
	}
	if cd.DependsOn != nil {
		container.DependsOn = cd.DependsOn
	} else {
		container.DependsOn = old.DependsOn
	}
	if cd.RepositoryCredentials != nil {
		container.RepositoryCredentials = cd.RepositoryCredentials
	} else {
		container.RepositoryCredentials = old.RepositoryCredentials
	}
	if cd.DockerSecurityOptions != nil {
		container.DockerSecurityOptions = cd.DockerSecurityOptions
	} else {
		container.DockerSecurityOptions = old.DockerSecurityOptions
	}
	if cd.DockerLabels != nil {
		container.DockerLabels = cd.DockerLabels
	} else {
		container.DockerLabels = old.DockerLabels
	}
	if cd.Essential != nil {
		container.Essential = cd.Essential
	} else {
		container.Essential = old.Essential
	}
	if cd.Cpu != nil {
		container.Cpu = cd.Cpu
	} else {
		container.Cpu = old.Cpu
	}
	if cd.Memory != nil {
		container.Memory = cd.Memory
	} else {
		container.Memory = old.Memory
	}
	if cd.MemoryReservation != nil {
		container.MemoryReservation = cd.MemoryReservation
	} else {
		container.MemoryReservation = old.MemoryReservation
	}
	if cd.ResourceRequirements != nil {
		container.ResourceRequirements = cd.ResourceRequirements
	} else {
		container.ResourceRequirements = old.ResourceRequirements
	}
	if cd.Ulimits != nil {
		container.Ulimits = cd.Ulimits
	} else {
		container.Ulimits = old.Ulimits
	}
	if cd.User != nil {
		container.User = cd.User
	} else {
		container.User = old.User
	}
	if cd.WorkingDirectory != nil {
		container.WorkingDirectory = cd.WorkingDirectory
	} else {
		container.WorkingDirectory = old.WorkingDirectory
	}
	if cd.Interactive != nil {
		container.Interactive = cd.Interactive
	} else {
		container.Interactive = old.Interactive
	}
	if cd.PseudoTerminal != nil {
		container.PseudoTerminal = cd.PseudoTerminal
	} else {
		container.PseudoTerminal = old.PseudoTerminal
	}
	if cd.ReadonlyRootFilesystem != nil {
		container.ReadonlyRootFilesystem = cd.ReadonlyRootFilesystem
	} else {
		container.ReadonlyRootFilesystem = old.ReadonlyRootFilesystem
	}
	if cd.Privileged != nil {
		container.Privileged = cd.Privileged
	} else {
		container.Privileged = old.Privileged
	}
	if cd.LinuxParameters != nil {
		container.LinuxParameters = cd.LinuxParameters
	} else {
		container.LinuxParameters = old.LinuxParameters
	}
	if cd.SystemControls != nil {
		container.SystemControls = cd.SystemControls
	} else {
		container.SystemControls = old.SystemControls
	}
	if cd.Hostname != nil {
		container.Hostname = cd.Hostname
	} else {
		container.Hostname = old.Hostname
	}
	if cd.ExtraHosts != nil {
		container.ExtraHosts = cd.ExtraHosts
	} else {
		container.ExtraHosts = old.ExtraHosts
	}
	if cd.MountPoints != nil {
		container.MountPoints = cd.MountPoints
	} else {
		container.MountPoints = old.MountPoints
	}
	if cd.VolumesFrom != nil {
		container.VolumesFrom = cd.VolumesFrom
	} else {
		container.VolumesFrom = old.VolumesFrom
	}
	if cd.Environment != nil {
		container.Environment = cd.Environment
	} else {
		container.Environment = old.Environment
	}
	if cd.Secrets != nil {
		container.Secrets = cd.Secrets
	} else {
		container.Secrets = old.Secrets
	}
	if cd.Links != nil {
		container.Links = cd.Links
	} else {
		container.Links = old.Links
	}
	if cd.PortMappings != nil {
		container.PortMappings = cd.PortMappings
	} else {
		container.PortMappings = old.PortMappings
	}
	if cd.DisableNetworking != nil {
		container.DisableNetworking = cd.DisableNetworking
	} else {
		container.DisableNetworking = old.DisableNetworking
	}
	if cd.DnsSearchDomains != nil {
		container.DnsSearchDomains = cd.DnsSearchDomains
	} else {
		container.DnsSearchDomains = old.DnsSearchDomains
	}
	if cd.DnsServers != nil {
		container.DnsServers = cd.DnsServers
	} else {
		container.DnsServers = old.DnsServers
	}
	if cd.HealthCheck != nil {
		container.HealthCheck = cd.HealthCheck
	} else {
		container.HealthCheck = old.HealthCheck
	}
	if cd.StartTimeout != nil {
		container.StartTimeout = cd.StartTimeout
	} else {
		container.StartTimeout = old.StartTimeout
	}
	if cd.StopTimeout != nil {
		container.StopTimeout = cd.StopTimeout
	} else {
		container.StopTimeout = old.StopTimeout
	}
	if cd.FirelensConfiguration != nil {
		container.FirelensConfiguration = cd.FirelensConfiguration
	} else {
		container.FirelensConfiguration = old.FirelensConfiguration
	}
	if cd.LogConfiguration != nil {
		container.LogConfiguration = cd.LogConfiguration
	} else {
		container.LogConfiguration = old.LogConfiguration
	}

	return container
}
