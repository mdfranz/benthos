package integration

import (
	"testing"
	"time"
)

var _ = registerIntegrationTest("nanomsg", func(t *testing.T) {
	t.Parallel()

	template := `
output:
  nanomsg:
    urls:
      - tcp://localhost:$PORT
    bind: false
    socket_type: $VAR1
    poll_timeout: 5s
    max_in_flight: $MAX_IN_FLIGHT

input:
  nanomsg:
    urls:
      - tcp://*:$PORT
    bind: true
    socket_type: $VAR2
    sub_filters: [ $VAR3 ]
`
	suite := integrationTests(
		integrationTestOpenClose(),
		integrationTestSendBatch(10),
		integrationTestStreamParallel(100),
	)
	suite.Run(
		t, template,
		testOptSleepAfterInput(500*time.Millisecond),
		testOptSleepAfterOutput(500*time.Millisecond),
		testOptVarOne("PUSH"),
		testOptVarTwo("PULL"),
	)
	t.Run("with max in flight", func(t *testing.T) {
		t.Parallel()
		suite.Run(
			t, template,
			testOptSleepAfterInput(500*time.Millisecond),
			testOptSleepAfterOutput(500*time.Millisecond),
			testOptVarOne("PUSH"),
			testOptVarTwo("PULL"),
			testOptMaxInFlight(10),
		)
	})
	t.Run("with pub sub", func(t *testing.T) {
		t.Parallel()
		suite.Run(
			t, template,
			testOptSleepAfterInput(500*time.Millisecond),
			testOptSleepAfterOutput(500*time.Millisecond),
			testOptVarOne("PUB"),
			testOptVarTwo("SUB"),
			testOptVarThree(`""`),
		)
	})
})
