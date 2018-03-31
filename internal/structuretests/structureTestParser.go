package structuretests

import (
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/container-structure-test/pkg/drivers"
	"github.com/GoogleCloudPlatform/container-structure-test/pkg/output"
	"github.com/GoogleCloudPlatform/container-structure-test/pkg/types/unversioned"
	"github.com/GoogleCloudPlatform/container-structure-test/pkg/types/v2"
	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"
	"github.com/zanetworker/dockument/pkg/labels"
)

//ParseTests is a method to read the tests from the docker image and parse it to match strcuture tests
func ParseTests(imageName, driver string, verboseOutput, supressOut bool) {
	if validateImageLocalAvailability(imageName) {

		driverImpl := drivers.InitDriverImpl("docker")
		if driverImpl == nil {
			log.Fatalf("Unsupported driver type: %s", driver)
		}

		args := &drivers.DriverConfig{
			Image: imageName,
		}
		structureTest := v2.StructureTest{}

		// add commandTests from docker image
		commandTests, _ := labels.GetImageCommandTests(imageName)
		if verboseOutput {
			pp.Println(commandTests)
		}

		for _, commandTest := range *commandTests {
			targetCommandTest := v2.CommandTest{}
			targetCommandTest.Name = commandTest.Name
			targetCommandTest.Command = commandTest.Command

			commandTestArgs := strings.Split(commandTest.Args[0], ",")
			if !containsEmptyStrings(commandTestArgs) {
				targetCommandTest.Args = commandTestArgs
			}

			excludedOutput := []string{commandTest.ExcludedOutput}
			if !containsEmptyStrings(excludedOutput) {
				targetCommandTest.ExcludedOutput = excludedOutput
			}

			excludedError := []string{commandTest.ExcludedError}
			if !containsEmptyStrings(excludedError) {
				targetCommandTest.ExcludedError = excludedError
			}

			expectedOuput := []string{commandTest.ExpectedOutput}
			if !containsEmptyStrings(expectedOuput) {
				targetCommandTest.ExpectedOutput = expectedOuput
			}

			expectedError := []string{commandTest.ExpectedError}
			if !containsEmptyStrings(expectedError) {
				targetCommandTest.ExpectedError = expectedError
			}

			err := targetCommandTest.Validate()
			if err != nil {
				log.Fatalf("invalid command structure Test, err: %s", err.Error())
			}
			structureTest.CommandTests = append(structureTest.CommandTests, targetCommandTest)
		}

		// add fileContentTests from docker image
		fileContentTests, _ := labels.GetImageFileContentTests(imageName)
		if verboseOutput {
			pp.Println(fileContentTests)
		}

		for _, fileContentTest := range *fileContentTests {
			targetFileContentTest := v2.FileContentTest{}

			targetFileContentTest.Name = fileContentTest.Name
			targetFileContentTest.Path = fileContentTest.Path

			expectedContents := []string{fileContentTest.ExpectedContents}
			if !containsEmptyStrings(expectedContents) {
				targetFileContentTest.ExpectedContents = []string{fileContentTest.ExpectedContents}
			}

			excludedContents := []string{fileContentTest.ExcludedContents}
			if !containsEmptyStrings(excludedContents) {
				targetFileContentTest.ExcludedContents = []string{fileContentTest.ExcludedContents}
			}

			err := targetFileContentTest.Validate()
			if err != nil {
				log.Fatalf("invalid File Content Test, err: %s", err.Error())
			}
			structureTest.FileContentTests = append(structureTest.FileContentTests, targetFileContentTest)
		}

		// add fileExistenceTests from docker image
		fileExistenceTests, _ := labels.GetImageFileExistenceTests(imageName)
		if verboseOutput {
			pp.Println(fileExistenceTests)
		}

		for _, fileExistenceTest := range *fileExistenceTests {
			targetFileExistenceTest := v2.FileExistenceTest{}

			targetFileExistenceTest.Name = fileExistenceTest.Name
			targetFileExistenceTest.Path = fileExistenceTest.Path

			targetFileExistenceTest.ShouldExist, _ = strconv.ParseBool(fileExistenceTest.ShouldExist)

			if !emptyString(fileExistenceTest.Permissions) {
				targetFileExistenceTest.Permissions = fileExistenceTest.Permissions
			}

			err := targetFileExistenceTest.Validate()
			if err != nil {
				log.Fatalf("invalid File Existence Test, err: %s", err.Error())
			}
			structureTest.FileExistenceTests = append(structureTest.FileExistenceTests, targetFileExistenceTest)
		}

		// get metadata tests from docker image
		metadataTests, _ := labels.GetImageMetadataTests(imageName)
		if verboseOutput {
			pp.Println(metadataTests)
		}

		for _, metadataTest := range *metadataTests {
			targetMetaDataTest := v2.MetadataTest{}
			targetMetaDataTest.Env = getEnvVar(metadataTest.Env)

			cmdArgs := strings.Split(metadataTest.Cmd, ",")
			if !containsEmptyStrings(cmdArgs) {
				targetMetaDataTest.Cmd = &cmdArgs
			}

			entrypointArgs := strings.Split(metadataTest.EntryPoint, ",")
			if !containsEmptyStrings(entrypointArgs) {
				targetMetaDataTest.Entrypoint = &entrypointArgs
			}

			exposedPorts := strings.Split(metadataTest.ExposedPorts, ",")
			if !containsEmptyStrings(exposedPorts) {
				targetMetaDataTest.ExposedPorts = exposedPorts
			}

			if !emptyString(metadataTest.Workdir) {
				targetMetaDataTest.Workdir = metadataTest.Workdir
			}

			volumes := strings.Split(metadataTest.Volumes, ",")
			if !containsEmptyStrings(volumes) {
				targetMetaDataTest.Volumes = volumes
			}

			structureTest.MetadataTest = targetMetaDataTest
		}

		structureTest.SetDriverImpl(driverImpl, *args)

		// execute tests
		results := structureTest.RunAll()
		out := &output.OutWriter{
			Verbose: verboseOutput,
			Quiet:   supressOut,
		}

		// execute metadata tests and add it to results
		metadataTestResult := structureTest.RunMetadataTests()
		results = append(results, metadataTestResult)

		// add all results
		fullResults := []*unversioned.FullResult{}
		fullResults = append(fullResults, &unversioned.FullResult{
			Results: results,
		})

		out.OutputResults(fullResults, imageName)

	} else {
		log.Error("Image is not available locally, trying to pull it from Dockerhub!")
		if err := pullDockerImage(imageName); err != nil {
			log.Fatalf("Failed to pull image from Dockerhub: %s", err.Error())
		}
	}

}
