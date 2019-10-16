package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type ContainerDefinition struct {
	Overwrite bool

	Name       string
	Image      *string
	EntryPoint *[]*string
	Command    *[]*string

	DependsOn             *[]*ecs.ContainerDependency
	RepositoryCredentials *ecs.RepositoryCredentials
	DockerSecurityOptions *[]*string
	DockerLabels          *map[string]*string

	Essential            *bool
	Cpu                  *int64
	Memory               *int64
	MemoryReservation    *int64
	ResourceRequirements *[]*ecs.ResourceRequirement
	Ulimits              *[]*ecs.Ulimit

	User                   *string
	WorkingDirectory       *string
	Interactive            *bool
	PseudoTerminal         *bool
	ReadonlyRootFilesystem *bool
	Privileged             *bool
	LinuxParameters        *ecs.LinuxParameters
	SystemControls         *[]*ecs.SystemControl

	Hostname    *string
	ExtraHosts  *[]*ecs.HostEntry
	MountPoints *[]*ecs.MountPoint
	VolumesFrom *[]*ecs.VolumeFrom

	Environment *[]*ecs.KeyValuePair
	Secrets     *[]*ecs.Secret

	Links        *[]*string
	PortMappings *[]*ecs.PortMapping

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

func (cd *ContainerDefinition) parse(containerDefinition *ecs.ContainerDefinition) *ContainerDefinition {
	if containerDefinition == nil {
		return cd
	}

	if cd.Image == nil {
		cd.Image = containerDefinition.Image
	}
	if cd.EntryPoint == nil {
		cd.EntryPoint = &containerDefinition.EntryPoint
	}
	if cd.Command == nil {
		cd.Command = &containerDefinition.Command
	}
	if cd.DependsOn == nil {
		cd.DependsOn = &containerDefinition.DependsOn
	}
	if cd.RepositoryCredentials == nil {
		cd.RepositoryCredentials = containerDefinition.RepositoryCredentials
	}
	if cd.DockerSecurityOptions == nil {
		cd.DockerSecurityOptions = &containerDefinition.DockerSecurityOptions
	}
	if cd.DockerLabels == nil {
		cd.DockerLabels = &containerDefinition.DockerLabels
	}
	if cd.Essential == nil {
		cd.Essential = containerDefinition.Essential
	}
	if cd.Cpu == nil {
		cd.Cpu = containerDefinition.Cpu
	}
	if cd.Memory == nil {
		cd.Memory = containerDefinition.Memory
	}
	if cd.MemoryReservation == nil {
		cd.MemoryReservation = containerDefinition.MemoryReservation
	}
	if cd.ResourceRequirements == nil {
		cd.ResourceRequirements = &containerDefinition.ResourceRequirements
	}
	if cd.Ulimits == nil {
		cd.Ulimits = &containerDefinition.Ulimits
	}
	if cd.User == nil {
		cd.User = containerDefinition.User
	}
	if cd.WorkingDirectory == nil {
		cd.WorkingDirectory = containerDefinition.WorkingDirectory
	}
	if cd.Interactive == nil {
		cd.Interactive = containerDefinition.Interactive
	}
	if cd.PseudoTerminal == nil {
		cd.PseudoTerminal = containerDefinition.PseudoTerminal
	}
	if cd.ReadonlyRootFilesystem == nil {
		cd.ReadonlyRootFilesystem = containerDefinition.ReadonlyRootFilesystem
	}
	if cd.Privileged == nil {
		cd.Privileged = containerDefinition.Privileged
	}
	if cd.LinuxParameters == nil {
		cd.LinuxParameters = containerDefinition.LinuxParameters
	}
	if cd.SystemControls == nil {
		cd.SystemControls = &containerDefinition.SystemControls
	}
	if cd.Hostname == nil {
		cd.Hostname = containerDefinition.Hostname
	}
	if cd.ExtraHosts == nil {
		cd.ExtraHosts = &containerDefinition.ExtraHosts
	}
	if cd.MountPoints == nil {
		cd.MountPoints = &containerDefinition.MountPoints
	}
	if cd.VolumesFrom == nil {
		cd.VolumesFrom = &containerDefinition.VolumesFrom
	}
	if cd.Environment == nil {
		cd.Environment = &containerDefinition.Environment
	}
	if cd.Secrets == nil {
		cd.Secrets = &containerDefinition.Secrets
	}
	if cd.Links == nil {
		cd.Links = &containerDefinition.Links
	}
	if cd.PortMappings == nil {
		cd.PortMappings = &containerDefinition.PortMappings
	}
	if cd.DisableNetworking == nil {
		cd.DisableNetworking = containerDefinition.DisableNetworking
	}
	if cd.DnsSearchDomains == nil {
		cd.DnsSearchDomains = containerDefinition.DnsSearchDomains
	}
	if cd.DnsServers == nil {
		cd.DnsServers = containerDefinition.DnsServers
	}
	if cd.HealthCheck == nil {
		cd.HealthCheck = containerDefinition.HealthCheck
	}
	if cd.StartTimeout == nil {
		cd.StartTimeout = containerDefinition.StartTimeout
	}
	if cd.StopTimeout == nil {
		cd.StopTimeout = containerDefinition.StopTimeout
	}
	if cd.FirelensConfiguration == nil {
		cd.FirelensConfiguration = containerDefinition.FirelensConfiguration
	}
	if cd.LogConfiguration == nil {
		cd.LogConfiguration = containerDefinition.LogConfiguration
	}

	return cd
}

func (cd *ContainerDefinition) unpack() *ecs.ContainerDefinition {
	containerDefinition := &ecs.ContainerDefinition{}

	containerDefinition.Name = &cd.Name

	if cd.Image != nil {
		containerDefinition.Image = cd.Image
	}
	if cd.EntryPoint != nil {
		containerDefinition.EntryPoint = *cd.EntryPoint
	}
	if cd.Command != nil {
		containerDefinition.Command = *cd.Command
	}
	if cd.DependsOn != nil {
		containerDefinition.DependsOn = *cd.DependsOn
	}
	if cd.RepositoryCredentials != nil {
		containerDefinition.RepositoryCredentials = cd.RepositoryCredentials
	}
	if cd.DockerSecurityOptions != nil {
		containerDefinition.DockerSecurityOptions = *cd.DockerSecurityOptions
	}
	if cd.DockerLabels != nil {
		containerDefinition.DockerLabels = *cd.DockerLabels
	}
	if cd.Essential != nil {
		containerDefinition.Essential = cd.Essential
	}
	if cd.Cpu != nil {
		containerDefinition.Cpu = cd.Cpu
	}
	if cd.Memory != nil {
		containerDefinition.Memory = cd.Memory
	}
	if cd.MemoryReservation != nil {
		containerDefinition.MemoryReservation = cd.MemoryReservation
	}
	if cd.ResourceRequirements != nil {
		containerDefinition.ResourceRequirements = *cd.ResourceRequirements
	}
	if cd.Ulimits != nil {
		containerDefinition.Ulimits = *cd.Ulimits
	}
	if cd.User != nil {
		containerDefinition.User = cd.User
	}
	if cd.WorkingDirectory != nil {
		containerDefinition.WorkingDirectory = cd.WorkingDirectory
	}
	if cd.Interactive != nil {
		containerDefinition.Interactive = cd.Interactive
	}
	if cd.PseudoTerminal != nil {
		containerDefinition.PseudoTerminal = cd.PseudoTerminal
	}
	if cd.ReadonlyRootFilesystem != nil {
		containerDefinition.ReadonlyRootFilesystem = cd.ReadonlyRootFilesystem
	}
	if cd.Privileged != nil {
		containerDefinition.Privileged = cd.Privileged
	}
	if cd.LinuxParameters != nil {
		containerDefinition.LinuxParameters = cd.LinuxParameters
	}
	if cd.SystemControls != nil {
		containerDefinition.SystemControls = *cd.SystemControls
	}
	if cd.Hostname != nil {
		containerDefinition.Hostname = cd.Hostname
	}
	if cd.ExtraHosts != nil {
		containerDefinition.ExtraHosts = *cd.ExtraHosts
	}
	if cd.MountPoints != nil {
		containerDefinition.MountPoints = *cd.MountPoints
	}
	if cd.VolumesFrom != nil {
		containerDefinition.VolumesFrom = *cd.VolumesFrom
	}
	if cd.Environment != nil {
		containerDefinition.Environment = *cd.Environment
	}
	if cd.Secrets != nil {
		containerDefinition.Secrets = *cd.Secrets
	}
	if cd.Links != nil {
		containerDefinition.Links = *cd.Links
	}
	if cd.PortMappings != nil {
		containerDefinition.PortMappings = *cd.PortMappings
	}
	if cd.DisableNetworking != nil {
		containerDefinition.DisableNetworking = cd.DisableNetworking
	}
	if cd.DnsSearchDomains != nil {
		containerDefinition.DnsSearchDomains = cd.DnsSearchDomains
	}
	if cd.DnsServers != nil {
		containerDefinition.DnsServers = cd.DnsServers
	}
	if cd.HealthCheck != nil {
		containerDefinition.HealthCheck = cd.HealthCheck
	}
	if cd.StartTimeout != nil {
		containerDefinition.StartTimeout = cd.StartTimeout
	}
	if cd.StopTimeout != nil {
		containerDefinition.StopTimeout = cd.StopTimeout
	}
	if cd.FirelensConfiguration != nil {
		containerDefinition.FirelensConfiguration = cd.FirelensConfiguration
	}
	if cd.LogConfiguration != nil {
		containerDefinition.LogConfiguration = cd.LogConfiguration
	}

	return containerDefinition
}
