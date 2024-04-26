package testutils

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestQueueProviderInstance struct {
	Container testcontainers.Container
	Context   context.Context
	Endpoint  string
}

func NewTestQueueProviderInstance() TestQueueProviderInstance {
	ctx := context.Background()

	sqsCustomConfPath := filepath.Join("../../sqs", "custom.conf")
	r, err := os.Open(sqsCustomConfPath)
	if err != nil {
		panic(err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "softwaremill/elasticmq",
		ExposedPorts: []string{"9324/tcp"},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            r,
				HostFilePath:      sqsCustomConfPath, // will be discarded internally
				ContainerFilePath: "/opt/elasticmq.conf",
				FileMode:          0o660,
			},
		},
		WaitingFor: wait.ForLog(`org\.elasticmq\.server\.Main\$ - === ElasticMQ server \(\d+\.\d+\.\d+\) started`).AsRegexp(),
	}

	elasticmqContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start elasticmqContainer: %s", err)
	}

	endpoint, err := elasticmqContainer.Endpoint(ctx, "http")
	if err != nil {
		log.Fatalf("Could not get elasticmqContainer endpoint: %s", err)
	}

	return TestQueueProviderInstance{Container: elasticmqContainer, Context: ctx, Endpoint: endpoint}
}

func (q *TestQueueProviderInstance) CleanUp() {
	// Clean up the container
	if err := q.Container.Terminate(q.Context); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}
