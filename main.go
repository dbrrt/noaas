// package main

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	r := gin.Default()
// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "pong",
// 		})
// 	})
// 	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
// }

package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/nomad/api"
)

func stringPtr(s string) *string {
	return &s
}

func intToPtr(i int) *int {
	return &i
}

func main() {
	// Create a new Nomad client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create Nomad client: %v", err)
	}

	j := createServiceJob()

	registerJob(client, j)

	// uri, err := getJobUri(client, "hello-world")

	// if err != nil {
	// 	log.Fatalf("Failed: %v", err)
	// }

	// fmt.Printf(uri)
}

// func createBatchJob() *api.Job {

// 	// Define a job
// 	job := &api.Job{
// 		ID:   stringPtr("example-batch"),
// 		Name: stringPtr("example-batch"),
// 		Type: stringPtr("batch"),
// 	}

// 	// Define task group
// 	taskGroup := &api.TaskGroup{
// 		Name:  stringPtr("example-group"),
// 		Count: intToPtr(1),
// 	}

// 	// // Define task
// 	task := &api.Task{
// 		Name:   "example-task",
// 		Driver: "docker",
// 		Config: map[string]interface{}{
// 			"image": "alpine",                                   // Docker image to run
// 			"args":  []string{"sh", "-c", "echo Hello, Nomad!"}, // Command
// 		},
// 		Resources: &api.Resources{
// 			CPU:      intToPtr(100),
// 			MemoryMB: intToPtr(128),
// 		},
// 	}

// 	// Add the task to the task group
// 	taskGroup.Tasks = []*api.Task{task}

// 	// Add the task group to the job
// 	job.TaskGroups = []*api.TaskGroup{taskGroup}

//		return job
//		// // Register the job with Nomad
//	}
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
				EmbeddedTmpl: stringPtr(`<h1>Hello, Nomad!</h1>
<ul>
  <li>Task: {{env "NOMAD_TASK_NAME"}}</li>
  <li>Group: {{env "NOMAD_GROUP_NAME"}}</li>
  <li>Job: {{env "NOMAD_JOB_NAME"}}</li>
  <li>Metadata value for foo: {{env "NOMAD_META_foo"}}</li>
  <li>Currently running on port: {{env "NOMAD_PORT_www"}}</li>
</ul>`),
				DestPath: stringPtr("local/index.html"), // Render template to the index.html file
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

func registerJob(client *api.Client, job *api.Job) {
	// Register the job with Nomad
	resp, _, err := client.Jobs().Register(job, nil)
	if err != nil {
		log.Fatalf("Failed to register job: %v", err)
	}

	// Output Job ID and Evaluate ID
	fmt.Printf("Job registered: ID=%s EvalID=%s\n", *job.ID, resp.EvalID)

	// Optional: Monitor job deployment status
	// waitForJobCompletion(client, *job.ID)
}

func getJobUri(client *api.Client, jobId string) (string, error) {
	// Get the job's allocation list stubs
	allocationStubs, _, err := client.Jobs().Allocations(jobId, false, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get allocations for job %s: %v", jobId, err)
	}

	// Iterate over allocation stubs to find the first running allocation
	var alloc *api.Allocation
	for _, allocStub := range allocationStubs {
		if allocStub.ClientStatus == "running" {
			// Fetch full allocation details using the allocation ID
			alloc, _, err = client.Allocations().Info(allocStub.ID, nil)
			if err != nil {
				return "", fmt.Errorf("failed to get allocation info for ID %s: %v", allocStub.ID, err)
			}
			break
		}
	}

	if alloc == nil {
		return "", fmt.Errorf("no running allocations found for job %s", jobId)
	}

	// Retrieve the node details for this allocation
	node, _, err := client.Nodes().Info(alloc.NodeID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get node info for node %s: %v", alloc.NodeID, err)
	}
	// Retrieve the IP address from the node's Resources or Network info
	var nodeIP string
	for _, address := range node.Resources.Networks {
		nodeIP = address.IP // Assuming there's an IPAddress field in the network struct
		break
	}

	if nodeIP == "" {
		return "", fmt.Errorf("no IP address found for node %s", alloc.NodeID)
	}
	// Get the first dynamically allocated port (if any) from the allocation's resources
	var port int
	for _, p := range alloc.AllocatedResources.Shared.Ports {
		port = p.Value
		break // Take the first available port
	}

	if port == 0 {
		return "", fmt.Errorf("no allocated ports found for allocation %s", alloc.ID)
	}

	// Construct the URI
	// Construct the URI
	uri := fmt.Sprintf("http://%s:%d", nodeIP, port)
	return uri, nil
}

// func waitForJobCompletion(client *api.Client, jobID string) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
// 	defer cancel()

// 	for {
// 		// Fetch job info to check status
// 		job, _, err := client.Jobs().Info(jobID, nil)
// 		if err != nil {
// 			log.Fatalf("Failed to fetch job info: %v", err)
// 		}

// 		// Check if job has completed
// 		if job.Status != nil && *job.Status == "complete" {
// 			fmt.Printf("Job %s completed successfully\n", jobID)
// 			return
// 		}

// 		// Check for other job statuses
// 		fmt.Printf("Job %s status: %s\n", jobID, *job.Status)
// 		time.Sleep(5 * time.Second) // Poll every 5 seconds

// 		// End the loop if the context expires
// 		if ctx.Err() != nil {
// 			log.Fatalf("Timed out waiting for job completion")
// 			return
// 		}
// 	}
// }
