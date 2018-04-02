

```go

fmt.Printf("Service List: %+v\n", services)

Service List: [{ID:yfhvks7d6pvdt8rcuzsto5xyp Meta:{Version:{Index:48} CreatedAt:2018-04-01 22:55:17.986857637 +0000 UTC UpdatedAt:2018-04-01 22:55:17.99087939 +0000 UTC} Spec:{Annotations:{Name:nist-mirror Labels:map[]} TaskTemplate:{ContainerSpec:0xc4203c8840 PluginSpec:<nil> Resources:0xc420398810 RestartPolicy:<nil> Placement:0xc420350370 Networks:[] LogDriver:<nil> ForceUpdate:0 Runtime:container} Mode:{Replicated:0xc42035e818 Global:<nil>} UpdateConfig:<nil> RollbackConfig:<nil> Networks:[] EndpointSpec:0xc420385dd0} PreviousSpec:<nil> Endpoint:{Spec:{Mode:vip Ports:[{Name: Protocol:tcp TargetPort:80 PublishedPort:8080 PublishMode:ingress}]} Ports:[{Name: Protocol:tcp TargetPort:80 PublishedPort:8080 PublishMode:ingress}] VirtualIPs:[{NetworkID:9tzbjb1r2lqzh465qdakmi9ab Addr:10.255.0.13/16}]} UpdateStatus:<nil>}]


RUN: docker service update nist-mirror --force
func (cli *Client) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error)

func (cli *Client) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error)

func (cli *Client) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error)


Service List: [{ID:yfhvks7d6pvdt8rcuzsto5xyp Meta:{Version:{Index:48} CreatedAt:2018-04-01 22:55:17.986857637 +0000 UTC UpdatedAt:2018-04-01 22:55:17.99087939 +0000 UTC}

Spec: {
    Annotations:{
        Name:nist-mirror Labels:map[]
    }
    TaskTemplate:{
        ContainerSpec:0xc4203c8840 PluginSpec:<nil> Resources:0xc420398810 RestartPolicy:<nil> Placement:0xc420350370 Networks:[] LogDriver:<nil> ForceUpdate:0 Runtime:container
    }
    Mode:{Replicated:0xc42035e818 Global:<nil>}
    UpdateConfig:<nil> RollbackConfig:<nil> Networks:[] EndpointSpec:0xc420385dd0} PreviousSpec:<nil> Endpoint:{
        Spec:{Mode:vip Ports:[{Name: Protocol:tcp TargetPort:80 PublishedPort:8080 PublishMode:ingress}]} Ports:[{Name: Protocol:tcp TargetPort:80 PublishedPort:8080 PublishMode:ingress}] VirtualIPs:[{NetworkID:9tzbjb1r2lqzh465qdakmi9ab Addr:10.255.0.13/16}]
    }
    UpdateStatus:<nil>}]

// taskSpec := swarm.TaskSpec{
// 	ForceUpdate: 1,
// 	// Runtime:     swarm.RuntimeContainer,
// 	Runtime: eachSvc.Spec.TaskTemplate.Runtime,
// }
// serviceSpec := swarm.ServiceSpec{
// 	TaskTemplate: taskSpec,
// }

```