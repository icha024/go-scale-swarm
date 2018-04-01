package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"io"
	"time"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	pf := fmt.Printf
	for _, container := range containers {
		pf("%s %s\n", container.ID[:10], container.Image)
	}

	evtOptions := types.EventsOptions{
		Since: time.Now().Format(time.RFC3339),
	}
	evtMsgChan, errChan := cli.Events(context.Background(), evtOptions)

	for {
		select {
		case err := <-errChan:
			if err != nil && err != io.EOF {
				fmt.Println(err)
			}
		case evt := <-evtMsgChan:
			pf("Received Event: Type=%s, Action=%s\n", evt.Type, evt.Action)

			if evt.Type != "abc" {
				pf("Rebalancing node...\n")

				services, svcListErr := cli.ServiceList(context.Background(), types.ServiceListOptions{})
				if svcListErr != nil {
					pf("Service List Error: %s\n", svcListErr)
					continue
				}

				for _, eachSvc := range services {
					pf("Found service %s (%s)\n", eachSvc.Spec.Annotations.Name, eachSvc.ID)

					svcID := eachSvc.ID
					swarmVersion := swarm.Version{
						Index: 19,
					}
					updateResp, updateErr := cli.ServiceUpdate(context.Background(), svcID, swarmVersion, swarm.ServiceSpec{}, types.ServiceUpdateOptions{})
					if updateErr != nil {
						pf("Service Update Error: %s\n", updateErr)
						continue
					}
					pf("Service Update Response: %s\n", updateResp)
				}
			}
		}
	}
}
