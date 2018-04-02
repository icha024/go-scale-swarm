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
	pf := fmt.Printf
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
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

			if evt.Type != "abc" { // FIXME: evt.Type == "node" && evt.Action == "create"
				pf("Rebalancing node...\n")

				services, svcListErr := cli.ServiceList(context.Background(), types.ServiceListOptions{})
				if svcListErr != nil {
					pf("Service List Error: %s\n", svcListErr)
					continue
				}

				for _, eachSvc := range services {
					pf("Found service %s (%s)\n", eachSvc.Spec.Annotations.Name, eachSvc.ID)

					// FIXME: Docker API error
					svcID := eachSvc.ID
					swarmVersion := swarm.Version{Index: 19}
					serviceSpec := eachSvc.Spec
					updateResp, updateErr := cli.ServiceUpdate(context.Background(), svcID, swarmVersion, serviceSpec, types.ServiceUpdateOptions{})
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
