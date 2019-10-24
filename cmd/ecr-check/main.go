package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/carash/ecs-deploy/ecr"
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
			Name:   "ecr-image",
			Usage:  "Full URL of the image",
			EnvVar: "PLUGIN_IMAGE",
		},
		cli.Int64Flag{
			Name:   "check-interval",
			Usage:  "Interval to check availability of image, defaults to 10 seconds",
			EnvVar: "PLUGIN_INTERVAL",
		},
		cli.Int64Flag{
			Name:   "check-timeout",
			Usage:  "Timeout when checking availability of image, defaults to 60 seconds",
			EnvVar: "PLUGIN_TIMEOUT",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	creds := ecr.Credential{}
	if c.IsSet("access-key") {
		s := c.String("access-key")
		creds.AWSAccessKeyID = &s
	}
	if c.IsSet("secret-key") {
		s := c.String("secret-key")
		creds.AWSSecretAccessKey = &s
	}
	if c.IsSet("assume-role-arn") {
		s := c.String("assume-role-arn")
		creds.AWSAssumeRoleARN = &s
	}
	if c.IsSet("aws-region") {
		s := c.String("aws-region")
		creds.AWSRegion = &s
	}

	image, err := parseECRImage(c.String("ecr-image"))
	if err != nil {
		return err
	}

	plugin := ecr.ImagePlugin{
		AWSCredential: creds,
		Image:         *image,
	}

	var interval, timeout int64
	if c.IsSet("check-interval") {
		interval = c.Int64("check-interval")
	} else {
		interval = 10
	}

	if c.IsSet("check-timeout") {
		timeout = c.Int64("check-timeout")
	} else {
		timeout = 60
	}

	return plugin.WaitForImage(interval, timeout)
}

var imageRegex, _ = regexp.Compile(`^\d{12}.dkr.ecr.[a-z]{2}-[a-z]+-\d{1,2}.amazonaws.com\/[\w-]+$`)
var taggedRegex, _ = regexp.Compile(`^\d{12}.dkr.ecr.[a-z]{2}-[a-z]+-\d{1,2}.amazonaws.com\/[\w-]+:[\w-]+$`)

func parseECRImage(image string) (*ecr.Image, error) {
	switch true {
	case imageRegex.MatchString(image):
		registry := image[:12]
		repository := strings.Split(image, "/")[1]
		return &ecr.Image{
			RegistryId:     &registry,
			RepositoryName: repository,
		}, nil
	case taggedRegex.MatchString(image):
		registry := image[:12]
		repotag := strings.Split(image, "/")[1]
		repository := strings.Split(repotag, ":")[0]
		tag := strings.Split(repotag, ":")[1]
		return &ecr.Image{
			RegistryId:     &registry,
			RepositoryName: repository,
			ImageTags:      &[]*string{&tag},
		}, nil
	}

	return nil, fmt.Errorf("Bad ECR Image Registry")
}
