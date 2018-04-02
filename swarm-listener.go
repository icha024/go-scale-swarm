package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	// "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"io"
	"os/exec"
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

			// if evt.Type != "abc" {
			if evt.Type == "node" && evt.Action == "create" {
				pf("Rebalancing node...\n")

				services, svcListErr := cli.ServiceList(context.Background(), types.ServiceListOptions{})
				if svcListErr != nil {
					pf("Service List Error: %s\n", svcListErr)
					continue
				}

				for _, eachSvc := range services {
					svcName := eachSvc.Spec.Annotations.Name
					pf("Found service %s (%s)\n", svcName, eachSvc.ID)

					// FIXME: Docker API error
					// svcID := eachSvc.ID
					// swarmVersion := swarm.Version{Index: 19}
					// serviceSpec := eachSvc.Spec
					// updateResp, updateErr := cli.ServiceUpdate(context.Background(), svcID, swarmVersion, serviceSpec, types.ServiceUpdateOptions{})
					// if updateErr != nil {
					// 	pf("Service Update Error: %s\n", updateErr)
					// 	continue
					// }
					// pf("Service Update Response: %s\n", updateResp)

					// Direct command with timeout
					cmd := exec.Command("docker", "service", "update", svcName, "--force")
					if cmdErr := cmd.Start(); err != nil {
						pf("%s service update error: %v\n", svcName, cmdErr)
						continue
					}

					// Wait for the process to finish or kill it after a timeout:
					cmdDone := make(chan error, 1)
					go func() {
						cmdDone <- cmd.Wait()
					}()
					select {
					case <-time.After(3 * time.Second):
						if killErr := cmd.Process.Kill(); err != nil {
							pf("Can not kill %s service update: %s\n", svcName, killErr)
						}
						pf("%s service update process killed as timeout reached\n", svcName)
					case cmdDoneErr := <-cmdDone:
						if cmdDoneErr != nil {
							pf("%s service update process finished with error = %v\n", svcName, cmdDoneErr)
						}
					}
				}
			}
		}
	}
}
