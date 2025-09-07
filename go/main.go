package main

import (
	"context"
	"fmt"
	"strings"

	containertypes "github.com/moby/moby/api/types/container"
	eventstypes "github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
		details, err := cli.ContainerInspect(ctx, container.ID)
		if err != nil {
			panic(err)
		}
		for _, env := range details.Config.Env {
			if strings.HasPrefix(env, "VIRTUAL_HOST=") {
				fmt.Println(env)
			}
		}
	}

	msgs, errs := cli.Events(ctx, eventstypes.ListOptions{})

	for {
		select {
		case err := <-errs:
			fmt.Println(err)
		case msg := <-msgs:
			if msg.Action == "start" {
				details, err := cli.ContainerInspect(ctx, msg.Actor.ID)
				if err != nil {
					panic(err)
				}
				for _, env := range details.Config.Env {
					if strings.HasPrefix(env, "VIRTUAL_HOST=") {
						fmt.Println(env)
					}
				}
			} else if msg.Action == "kill" || msg.Action == "stop" {
				details, err := cli.ContainerInspect(ctx, msg.Actor.ID)
				if err != nil {
					panic(err)
				}
				for _, env := range details.Config.Env {
					if strings.HasPrefix(env, "VIRTUAL_HOST=") {
						fmt.Println(env)
					}
				}
			}
		}
	}
}
