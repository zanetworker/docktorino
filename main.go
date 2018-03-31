package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/zanetworker/docktorino/internal/structuretests"
)

func main() {
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
				panic(err) // TODO(ahmetb) handle gracefully
			}

		case e := <-ch:
			// tag will be in format IMAGE:TAG or IMAGE:latest as it comes
			// from the Docker API (v1.32 at the time of writing).
			tag := e.Actor.Attributes["name"]
			structuretests.ParseTests(tag, "docker", false, false)
			// tagCh <- tag
		case <-ctx.Done():
			break
		}
	}
}
