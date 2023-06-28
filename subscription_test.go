// @@@SNIPSTART subscription-workflow-go-subscribe-test
package subscribeemails

import (
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func Test_CanceledSubscriptionWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := EmailDetails{
		EmailAddress:      "example@temporal.io",
		Message:           "This is a test to see if the Workflow cancels. This is dependent on the bool variable in the testDetails struct.",
		IsSubscribed:      false,
		SubscriptionCount: 4,
	}

	env.RegisterWorkflow(SubscriptionWorkflow)
	env.RegisterActivity(SendEmail)

	// Execute and cancel Workflow
	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
	env.CancelWorkflow()

}
// @@@SNIPEND
