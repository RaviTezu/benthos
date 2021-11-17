package pulsar_test

import (
	"fmt"
	"testing"
	"time"

	bpulsar "github.com/Jeffail/benthos/v3/internal/impl/pulsar"
	"github.com/Jeffail/benthos/v3/internal/integration"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPulsar(t *testing.T) {
	integration.CheckSkip(t)
	t.Parallel()

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	pool.MaxWait = time.Second * 30
	if dline, exists := t.Deadline(); exists {
		pool.MaxWait = time.Until(dline) / 2
	}

	resource, err := pool.Run("apachepulsar/pulsar-standalone", "latest", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, pool.Purge(resource))
	})

	resource.Expire(900)
	require.NoError(t, pool.Retry(func() error {
		client, err := pulsar.NewClient(pulsar.ClientOptions{
			URL:    fmt.Sprintf("pulsar://localhost:%v/", resource.GetPort("6650/tcp")),
			Logger: bpulsar.NoopLogger(),
		})
		if err != nil {
			return err
		}
		prod, err := client.CreateProducer(pulsar.ProducerOptions{
			Topic: "benthos-connection-test",
		})
		if err == nil {
			prod.Close()
		}
		client.Close()
		return err
	}))

	template := `
output:
  pulsar:
    url: pulsar://localhost:$PORT/
    topic: "topic-$ID"
    max_in_flight: $MAX_IN_FLIGHT

input:
  pulsar:
    url: pulsar://localhost:$PORT/
    topics: [ "topic-$ID" ]
    subscription_name: "sub-$ID"
`
	suite := integration.StreamTests(
		integration.StreamTestOpenClose(),
		integration.StreamTestSendBatch(10),
		integration.StreamTestStreamSequential(1000),
		integration.StreamTestStreamParallel(1000),
		integration.StreamTestStreamParallelLossy(1000),
		integration.StreamTestStreamParallelLossyThroughReconnect(1000),
		integration.StreamTestAtLeastOnceDelivery(),
	)
	suite.Run(
		t, template,
		integration.StreamTestOptSleepAfterInput(500*time.Millisecond),
		integration.StreamTestOptSleepAfterOutput(500*time.Millisecond),
		integration.StreamTestOptPort(resource.GetPort("6650/tcp")),
	)
	t.Run("with max in flight", func(t *testing.T) {
		t.Parallel()
		suite.Run(
			t, template,
			integration.StreamTestOptSleepAfterInput(500*time.Millisecond),
			integration.StreamTestOptSleepAfterOutput(500*time.Millisecond),
			integration.StreamTestOptPort(resource.GetPort("6650/tcp")),
			integration.StreamTestOptMaxInFlight(10),
		)
	})
}