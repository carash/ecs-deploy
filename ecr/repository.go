package ecr

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecr"
)

type Image struct {
	RegistryId     *string
	RepositoryName string
	ImageDigest    *string
	ImageTags      *[]*string
}

func (i *Image) DockerTag() string {
	repo := i.RepositoryName
	if i.RegistryId != nil {
		repo = fmt.Sprintf("%s.dkr.ecr.region.amazonaws.com/%s", *i.RegistryId, i.RepositoryName)
	}

	if i.ImageTags != nil && len(*i.ImageTags) == 1 {
		return fmt.Sprintf("%s:%s", repo, *(*i.ImageTags)[0])
	}

	return i.RepositoryName
}

func (i *Image) isValid() error {
	if i.RepositoryName == "" {
		return fmt.Errorf("Image must have a repository")
	}

	return nil
}

func (i *Image) Find(reg *ecr.ECR) ([]*ecr.ImageDetail, error) {
	if err := i.isValid(); err != nil {
		return nil, err
	}

	input := i.unpackDescribeInput()
	imgout, err := reg.DescribeImages(input)
	if err != nil {
		return nil, err
	}

	return imgout.ImageDetails, nil
}

func (i *Image) unpackDescribeInput() *ecr.DescribeImagesInput {
	describeImagesInput := &ecr.DescribeImagesInput{}

	describeImagesInput.RepositoryName = &i.RepositoryName
	describeImagesInput.RegistryId = i.RegistryId
	if i.ImageDigest != nil {
		describeImagesInput.ImageIds = append(describeImagesInput.ImageIds, &ecr.ImageIdentifier{ImageDigest: i.ImageDigest})
	}
	if i.ImageTags != nil {
		for _, tag := range *i.ImageTags {
			describeImagesInput.ImageIds = append(describeImagesInput.ImageIds, &ecr.ImageIdentifier{ImageTag: tag})
		}
	}

	return describeImagesInput
}
