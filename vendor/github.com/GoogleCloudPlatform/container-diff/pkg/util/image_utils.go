/*
Copyright 2017 Google, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"io/ioutil"
	"regexp"

	"github.com/sirupsen/logrus"
)

const LatestTag string = ":latest"

func GetImageLayers(pathToImage string) []string {
	layers := []string{}
	contents, err := ioutil.ReadDir(pathToImage)
	if err != nil {
		logrus.Error(err.Error())
	}

	for _, file := range contents {
		if file.IsDir() {
			layers = append(layers, file.Name())
		}
	}
	return layers
}

// checks to see if an image string contains a tag.
func HasTag(image string) bool {
	tagRegex := regexp.MustCompile(".*:[^/]+$")
	return tagRegex.MatchString(image)
}
