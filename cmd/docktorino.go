// Copyright © 2018 Adek Zaalouk <zanetworker@adelzaalouk.me>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/zanetworker/docktorino/internal/environment"
)

var (
	settings environment.EnvSettings
)

const globalUsage = `
Docktorino is tool to help you continiously test your Docker images while building them locally.
`

const (
	//FILE is used to tell dockument to fetch labels from a Dockerfile
	FILE = "file"
	//IMAGE is used to tell dockument to fetch labels from a Docker image
	IMAGE = "image"
)

var docktorinoLogo = `
_____             _    _             _             
|  __ \           | |  | |           (_)            
| |  | | ___   ___| | _| |_ ___  _ __ _ _ __   ___  
| |  | |/ _ \ / __| |/ / __/ _ \| '__| | '_ \ / _ \ 
| |__| | (_) | (__|   <| || (_) | |  | | | | | (_) |
|_____/ \___/ \___|_|\_\\__\___/|_|  |_|_| |_|\___/ 
                                                                                                              
`

func newRootCmd(args []string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	docktorinoCmd := &cobra.Command{
		Use:   "docktorino",
		Short: "Docktorino helps you contioniously test your docker images during local developement",
		Long:  globalUsage,
		Run:   runDocktorino,
	}

	flags := docktorinoCmd.PersistentFlags()
	settings.AddFlags(flags)
	out := docktorinoCmd.OutOrStdout()

	docktorinoCmd.AddCommand(
		newStartCmd(out),
	)

	// set defaults from environment
	settings.Init(flags)

	return docktorinoCmd
}

func printLogo(logoToPrint string) {
	figletColoured := aec.GreenF.Apply(logoToPrint)
	if runtime.GOOS == "windows" {
		figletColoured = aec.BlueF.Apply(logoToPrint)
	}
	if _, err := fmt.Println(figletColoured); err != nil {
		log.Fatalf("Failed to print Docktorino figlet, error: %s", err.Error())
	}
}

func returnWithError(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func runDocktorino(cmd *cobra.Command, args []string) {
	printLogo(docktorinoLogo)
	if len(args) == 0 {
		cmd.Help()
	}
}

//Execute command for dockument CLI
func main() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	err := doc.GenMarkdownTree(cmd, "./doc")
	if err != nil {
		log.Fatal(err)

	}
}
