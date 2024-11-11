package nomad

import (
	"dbrrt/noaas/readuri"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/nomad/api"
)

func CreateAJobAndGetUri(jobNameParam string, uriParam string, script bool) (string, error) {
	// Create a new Nomad client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return "", fmt.Errorf("failed to create Nomad client: %v", err)
	}

	// Create and register the job
	job, err := createServiceJob(jobNameParam, uriParam, script)
	if err != nil {
		return "", fmt.Errorf("error during service job creation: %v", err)
	}

	allocID, err := registerJobAndGetAllocationID(client, job)
	if err != nil {
		return "", fmt.Errorf("failed to register job / retrieve allocation ID: %v", err)
	}

	// Get allocation info using the allocation ID
	allocation, _, err := client.Allocations().Info(allocID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve allocation info: %v", err)
	}

	// Find the URI for the "www" dynamic port
	var uri string
	if allocation.AllocatedResources != nil {
		for _, network := range allocation.AllocatedResources.Shared.Networks {
			for _, port := range network.DynamicPorts {
				if port.Label == "www" { // Look for the "www" port label
					// Ensure IP is available
					if network.IP != "" {
						uri = fmt.Sprintf("%s:%d", network.IP, port.Value)
						break
					} else {
						fmt.Println("IP not yet available; waiting...")
						time.Sleep(5 * time.Second)
					}
				}
			}
			// Exit the outer loop if URI has been set
			if uri != "" {
				break
			}
		}
	}

	if uri == "" {
		return "", fmt.Errorf("no URI found for service www")
	} else {
		return uri, nil
	}
}

func registerJobAndGetAllocationID(client *api.Client, job *api.Job) (string, error) {
	// Register the job with Nomad
	resp, _, err := client.Jobs().Register(job, nil)
	if err != nil {
		return "", fmt.Errorf("failed to register job: %v", err)
	}

	// Output Job ID and Evaluation ID
	fmt.Printf("Job registered: ID=%s EvalID=%s\n", *job.ID, resp.EvalID)

	var allocID string

	// Poll until at least one allocation is available
	for {
		// Fetch allocations for the job
		allocs, _, err := client.Jobs().Allocations(*job.ID, false, nil)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve allocations: %v", err)
		}

		// Check if allocations are available
		if len(allocs) > 0 {
			allocID = allocs[0].ID // Get the ID of the first allocation
			fmt.Printf("Allocation ID found: %s\n", allocID)
			break
		}

		// Wait before polling again
		fmt.Println("Waiting for allocation to be created...")
		time.Sleep(5 * time.Second)
	}

	return allocID, nil
}

func createServiceJob(jobName string, uri string, script bool) (*api.Job, error) {
	// Define the service job
	job := &api.Job{
		ID:          stringPtr(uuid.New().String()),
		Name:        stringPtr(jobName),
		Type:        stringPtr("service"),
		Datacenters: []string{"*"}, // Specifies that this job can run in any datacenter
		Meta: map[string]string{
			"CreatedAt": time.Now().String(), // Meta information
		},
	}

	// Define task group
	taskGroup := &api.TaskGroup{
		Name:  stringPtr("servers"),
		Count: intToPtr(1), // Number of instances of this group to run
	}

	// Define the network configuration with dynamic port
	network := &api.NetworkResource{
		DynamicPorts: []api.Port{
			{
				Label: "www", // Label for the network port
			},
		},
	}

	// Define the service to expose the port
	service := &api.Service{
		Provider:  "nomad",
		PortLabel: "www", // Use the dynamic port label for the service
	}

	payloadRemote, err := readuri.ReadRemoteUriPayload(uri, script)

	if err != nil {
		return nil, fmt.Errorf("error during remote read payload")
	}

	// Define task
	task := &api.Task{
		Name:   "web",
		Driver: "docker",
		Config: map[string]interface{}{
			"image":   "busybox:1",                                                     // Docker image to use
			"command": "httpd",                                                         // Command to start the web server
			"args":    []string{"-v", "-f", "-p", "${NOMAD_PORT_www}", "-h", "/local"}, // Use dynamic port
			"ports":   []string{"www"},                                                 // Reference the dynamic port
		},
		Templates: []*api.Template{
			{
				EmbeddedTmpl: stringPtr(payloadRemote),
				DestPath:     stringPtr("local/index.html"), // Render template to the index.html file
			},
		},
		Resources: &api.Resources{
			CPU:      intToPtr(50), // CPU allocation in Mhz
			MemoryMB: intToPtr(64), // Memory allocation
		},
	}

	// Attach the network and service configurations to the task group
	taskGroup.Networks = []*api.NetworkResource{network}
	taskGroup.Services = []*api.Service{service}
	taskGroup.Tasks = []*api.Task{task}

	// Add the task group to the job
	job.TaskGroups = []*api.TaskGroup{taskGroup}

	return job, nil
}
