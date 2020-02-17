package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	cred "github.com/carash/ecs-deploy/credential"
	"github.com/carash/ecs-deploy/ecs"
	"github.com/urfave/cli"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "AWS ECS Deploy"
	app.Usage = "AWS ECS Deploy"
	app.Action = run
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "AWS access key",
			EnvVar: "PLUGIN_ACCESS_KEY,ECS_ACCESS_KEY,AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "AWS secret key",
			EnvVar: "PLUGIN_SECRET_KEY,ECS_SECRET_KEY,AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "assume-role-arn",
			Usage:  "AWS secret key",
			EnvVar: "PLUGIN_ASSUME_ROLE_ARN",
		},
		cli.StringFlag{
			Name:   "aws-region",
			Usage:  "aws region",
			EnvVar: "PLUGIN_AWS_REGION,AWS_DEFAULT_REGION",
		},
		cli.StringFlag{
			Name:   "cluster",
			Usage:  "AWS ECS cluster",
			EnvVar: "PLUGIN_CLUSTER",
		},
		cli.StringFlag{
			Name:   "service",
			Usage:  "Service to act on",
			EnvVar: "PLUGIN_SERVICE",
		},
		cli.Int64Flag{
			Name:   "desired-count",
			Usage:  "The number of instantiations of the specified task definition to place and keep running on your cluster",
			EnvVar: "PLUGIN_DESIRED_COUNT",
		},
		cli.StringSliceFlag{
			Name:   "deployment-configuration",
			Usage:  "Deployment parameters that control how many tasks run during the deployment and the ordering of stopping and starting tasks",
			EnvVar: "PLUGIN_DEPLOYMENT_CONFIGURATION",
		},
		cli.IntFlag{
			Name:   "health-check-grace-period",
			Usage:  "Number of seconds to hold off health checks",
			EnvVar: "PLUGIN_HEALTH_CHECK_GRACE_PREIOD",
		},
		cli.StringFlag{
			Name:   "container-name",
			Usage:  "Container name",
			EnvVar: "PLUGIN_CONTAINER",
		},
		cli.StringFlag{
			Name:   "docker-image",
			Usage:  "image to use",
			EnvVar: "PLUGIN_IMAGE",
		},
		cli.Int64Flag{
			Name:   "timeout",
			Usage:  "Timeout to wait for healthy check",
			EnvVar: "PLUGIN_TIMEOUT",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	creds := cred.Credential{}
	creds.AWSAccessKeyID = c.String("access-key")
	creds.AWSSecretAccessKey = c.String("secret-key")
	creds.AWSAssumeRoleARN = c.String("assume-role-arn")
	creds.AWSRegion = c.String("aws-region")

	service := ecs.Service{Service: c.String("service")}
	if c.IsSet("cluster") {
		s := c.String("cluster")
		service.Cluster = &s
	}
	if c.IsSet("desired-count") {
		i := c.Int64("desired-count")
		service.DesiredCount = &i
	}
	if c.IsSet("deployment-configuration") {
		dc := awsecs.DeploymentConfiguration{}
		for _, s := range c.StringSlice("deployment-configuration") {
			if ok, _ := regexp.MatchString(`minimumHealthyPercent=\d+`, s); ok {
				p, _ := strconv.ParseInt(strings.Split(s, "=")[1], 10, 64)
				dc.MinimumHealthyPercent = &p
			} else if ok, _ := regexp.MatchString(`maximumPercent=\d+`, s); ok {
				p, _ := strconv.ParseInt(strings.Split(s, "=")[1], 10, 64)
				dc.MaximumPercent = &p
			}
		}

		service.DeploymentConfiguration = &dc
	}
	if c.IsSet("health-check-grace-period") {
		i := c.Int64("health-check-grace-period")
		service.HealthCheckGracePeriodSeconds = &i
	}

	if c.IsSet("container-name") || c.IsSet("docker-image") {
		task := ecs.TaskDefinition{}

		container := ecs.ContainerDefinition{Name: c.String("container-name")}
		if c.IsSet("docker-image") {
			s := c.String("docker-image")
			container.Image = &s
		}

		task.ContainerDefinitions = []*ecs.ContainerDefinition{&container}
		service.TaskDefinition = &task
	}

	plugin := ecs.ServicePlugin{
		AWSCredential: creds,
		Service:       service,
	}

	var timeout int64
	if c.IsSet("timeout") {
		timeout = c.Int64("timeout")
	} else {
		timeout = 600
	}

	return plugin.UpdateService(timeout)
}
