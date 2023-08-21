// @@@SNIPSTART subscription-workflow-go-subscribe-test
package subscribeemails

import (
	"testing"
	"time"

	"go.temporal.io/sdk/testsuite"
)

func Test_CanceledSubscriptionWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := EmailDetails{
		EmailAddress:      "example@temporal.io",
		Message:           "This is a test to see if the Workflow cancels. This is dependent on the bool variable in the testDetails struct.",
		IsSubscribed:      true,
		SubscriptionCount: 12,
	}

	// set delayed callback to allow time for cancellation.
	env.RegisterDelayedCallback(func() {
		env.CancelWorkflow()
	}, 5 * time.Second)

	env.RegisterWorkflow(SubscriptionWorkflow)
	env.RegisterActivity(SendEmail)

	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
}




// @@@SNIPEND
