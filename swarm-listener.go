package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
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

			// if evt.Type != "node" && evt.Action != "create" {
			// 	continue
			// }
			pf("Rebalancing node...\n")

			services, svcListErr := cli.ServiceList(context.Background(), types.ServiceListOptions{})
			if svcListErr != nil {
				pf("Service List Error: %s\n", svcListErr)
				continue
			}

			for _, eachSvc := range services {
				svcName := eachSvc.Spec.Annotations.Name
				pf("Found service %s (%s)\n", svcName, eachSvc.ID)

				// No need to update for global.
				replicatedMode := eachSvc.Spec.Mode.Replicated
				if replicatedMode == nil {
					continue
				}

				// Direct command with timeout. FIXME: Use Docker API.
				cmd := exec.Command("docker", "service", "update", svcName, "--force")
				if cmdErr := cmd.Start(); err != nil {
					pf("%s service update error: %v\n", svcName, cmdErr)
					continue
				}

				cmdDone := make(chan error, 1)
				go func() {
					cmdDone <- cmd.Wait()
				}()
				select {
				case <-time.After(3 * time.Second): // FIXME 60 sec.
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
