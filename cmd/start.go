package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zanetworker/docktorino/internal/structuretests"
)

var startCmdDesc = `
This command starts to the docktorino listener which triggers a tests each time the image of choice is built`

type startCmd struct {
	image   string
	verbose bool
}

func newStartCmd(out io.Writer) *cobra.Command {
	startCmdParams := &startCmd{}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "starts a listener that triggers tests for the target image of choice",
		Long:  startCmdDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return startCmdParams.run()
		},
	}

	f := startCmd.Flags()

	f.StringVarP(&startCmdParams.image, "image", "i", "", "the image you wish to trigger tests for")
	f.BoolVarP(&startCmdParams.verbose, "verbose", "v", false, "verbose testing output")

	return startCmd
}

func (s *startCmd) run() error {
	if len(s.image) != 0 {
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			sig := <-signalChan
			log.Printf("received signal: %s", sig.String())
			cancel()
		}()

		tagEvent := filters.NewArgs()
		tagEvent.Add("type", "image")
		tagEvent.Add("event", "tag")

		d, err := dockerclient.NewEnvClient()
		if err != nil {
			log.Fatal(errors.Wrap(err, "cannot create docker client"))
		}
		_, err = d.ServerVersion(ctx)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to connect to docker api"))
		}

		ch, errCh := d.Events(ctx, types.EventsOptions{Filters: tagEvent})
		for {
			select {
			case err := <-errCh:
				select {
				case <-ctx.Done():
					log.Println("stopping event listener due to cancellation")
					os.Exit(0)
				default:
					panic(err)
				}

			case e := <-ch:
				tag := e.Actor.Attributes["name"]
				if tag == s.image {
					structuretests.ParseTests(tag, "docker", s.verbose, false)
				}
			case <-ctx.Done():
				break
			}
		}
	} else {
		log.Println("Please provide a valid image name e.g., zanetworker/dockument:latest")
	}
	return nil
}
