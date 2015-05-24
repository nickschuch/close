package ecs_backend

import (
  "gopkg.in/alecthomas/kingpin.v1"
  "github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ecs"
  log "github.com/Sirupsen/logrus"

  "github.com/nickschuch/close/backend"
)

var (
	cliECSRegion  = kingpin.Flag("ecs-region", "The region to connect to.").Default("ap-southeast-2").OverrideDefaultFromEnvar("ECS_REGION").String()
	cliECSCluster = kingpin.Flag("ecs-cluster", "The cluster to get the list from.").Default("default").OverrideDefaultFromEnvar("ECS_CLUSTER").String()
)

type ECSBackend struct {}

func init() {
	backend.Register("ecs", &ECSBackend{})
}

func (o *ECSBackend) List(id string) map[string]string {
  var list map[string]string

  client := getECSClient()

  // Get the list of task ID's which we can use to query for all the information.
  tasksInput := &ecs.ListTasksInput{}
	tasks, err := client.ListTasks(tasksInput)
  if err != nil {
		log.Fatalf("Cannot get a list of containers from the ECS backend.")
	}

  // Get all the tasks information which we can use to extract the environment
  // variables for later.
  describeInput := &ecs.DescribeTasksInput{
		Cluster: aws.String(*cliECSCluster),
		Tasks:   tasks.TaskARNs,
	}
	described, err := client.DescribeTasks(describeInput)
  if err != nil {
		log.Fatalf("Cannot get a list of containers from the ECS backend.")
	}

  // Loop over the containers and build a list of urls to hit.
  for _, t := range described.Tasks {
    for _, c := range t.Containers {
      // Ensure this container has the required environment variable to be
      // exposed through the load balancer.
      url := getContainerEnv(*t.TaskDefinitionARN, *c.Name, id)
      if url == "" {
        continue
      }

      // Add it to the list.
      list[*t.TaskARN] = url
    }
  }

  return list
}

func (o *ECSBackend) Close(id string) {
  client := getECSClient()

  // We now create a request and pass it off to ECS so we can remove a task.
	// https://github.com/awslabs/aws-sdk-go/blob/master/service/ecs/examples_test.go#L617
	params := &ecs.StopTaskInput{
		Cluster: aws.String(*cliECSCluster),
		Task:    aws.String(id),
	}
	_, err := client.StopTask(params)
  if err != nil {
		log.Fatalf("Could not stop the container: %v", id)
	}
}

func getECSClient() *ecs.ECS {
	client := ecs.New(&aws.Config{Region: *cliECSRegion})
	return client
}

func getContainerEnv(definition string, name string, key string) string {
	client := getECSClient()

	tasksDefInput := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(definition),
	}
	tasksDefOutput, err := client.DescribeTaskDefinition(tasksDefInput)
  if err != nil {
		log.Fatalf("Could not find the task definition.")
	}

	for _, c := range tasksDefOutput.TaskDefinition.ContainerDefinitions {
		if *c.Name != name {
			continue
		}

		// Now we know we can look for the environment variable.
		for _, e := range c.Environment {
			if *e.Name == key {
				return *e.Value
			}
		}
	}

	return ""
}
