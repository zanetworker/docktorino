package structuretests

import (
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/container-structure-test/pkg/types/unversioned"
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

func pullDockerImage(imageName string) error {
	d, err := docker.NewClientFromEnv()

	imageParts := strings.Split(imageName, ":")
	repo, tag := imageParts[0], imageParts[1]
	err = d.PullImage(docker.PullImageOptions{
		Repository:   repo,
		Tag:          tag,
		OutputStream: os.Stdout,
	}, docker.AuthConfiguration{})

	if err != nil {
		return err
	}
	return nil
}

func validateImageLocalAvailability(imageName string) bool {
	d, err := docker.NewClientFromEnv()
	images, err := d.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		log.Fatalf("Failed to list images %s", err.Error())
	}

	for _, image := range images {
		if containsImageName(image.RepoTags, imageName) {
			return true
		}
	}
	return false
}

// Checks if RepoTags of an image contains the requested image tag
func containsImageName(imageRepoTags []string, requestedImage string) bool {
	for _, tag := range imageRepoTags {
		nameWithoutTag := strings.Split(tag, ":")[0]
		if tag == requestedImage || nameWithoutTag == requestedImage {
			return true
		}
	}
	return false
}

func containsEmptyStrings(sliceToCheck []string) bool {
	for _, value := range sliceToCheck {
		if len(value) == 0 {
			return true
		}
	}
	return false
}

func emptyString(stringToCheck string) bool {
	if len(stringToCheck) == 0 {
		return true
	}
	return false
}

func getEnvVar(envVarsToConvert string) []unversioned.EnvVar {
	envVars := strings.Split(envVarsToConvert, ",")
	envVarsToReturn := []unversioned.EnvVar{}
	for _, env := range envVars {
		envVarToReturn := unversioned.EnvVar{}
		keyAndValue := strings.Split(env, "=")
		envVarToReturn.Key = keyAndValue[0]
		envVarToReturn.Value = keyAndValue[1]
		envVarsToReturn = append(envVarsToReturn, envVarToReturn)
	}
	return envVarsToReturn
}
