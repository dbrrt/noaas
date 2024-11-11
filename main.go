package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/nomad/api"
)

func main() {
	// Create a new Nomad client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create Nomad client: %v", err)
	}

	// Create and register the job
	job := createServiceJob()
	allocID, err := registerJobAndGetAllocationID(client, job)
	if err != nil {
		log.Fatalf("Failed to get allocation ID: %v", err)
	}

	// Get allocation info using the allocation ID
	allocation, _, err := client.Allocations().Info(allocID, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve allocation info: %v", err)
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
		fmt.Println("No URI found for the 'www' port.")
	} else {
		fmt.Printf("Service available at URI: %s\n", uri)
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

func createServiceJob() *api.Job {
	// Define the service job
	job := &api.Job{
		ID:          stringPtr("hello-world"),
		Name:        stringPtr("hello-world"),
		Type:        stringPtr("service"),
		Datacenters: []string{"*"}, // Specifies that this job can run in any datacenter
		Meta: map[string]string{
			"foo": "bar", // Meta information
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
				EmbeddedTmpl: stringPtr(`<h1>Hello, Nomad!</h1>`),
				DestPath:     stringPtr("local/index.html"), // Render template to the index.html file
			},
		},
		Resources: &api.Resources{
			CPU:      intToPtr(50), // CPU allocation
			MemoryMB: intToPtr(64), // Memory allocation
		},
	}

	// Attach the network and service configurations to the task group
	taskGroup.Networks = []*api.NetworkResource{network}
	taskGroup.Services = []*api.Service{service}
	taskGroup.Tasks = []*api.Task{task}

	// Add the task group to the job
	job.TaskGroups = []*api.TaskGroup{taskGroup}

	return job
}

func stringPtr(s string) *string {
	return &s
}

func intToPtr(i int) *int {
	return &i
}
