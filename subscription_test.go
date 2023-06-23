// @@@SNIPSTART subscription-workflow-go-subscribe-test
package subscribeemails

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_SuccessfulSubscriptionWorkflow (t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := EmailDetails{
		EmailAddress: "example@temporal.io",
		Message: "This is a test to see if the Workflow successfully runs.",
		IsSubscribed: true,
		SubscriptionCount: 0,
		MaxSubscriptionPeriods: 12,
	}

	env.RegisterWorkflow(SubscriptionWorkflow)
	env.RegisterActivity(SendEmail)
	
	// Execute Workflow
	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
	require.NoError(t, env.GetWorkflowError())
}

func Test_CanceledSubscriptionWorkflow (t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := EmailDetails{
		EmailAddress: "example@temporal.io",
		Message: "This is a test to see if the Workflow cancels. This is dependent on the bool variable in the testDetails struct.",
		IsSubscribed: false,
		SubscriptionCount: 4,
		MaxSubscriptionPeriods: 12,
	}

	env.RegisterWorkflow(SubscriptionWorkflow)
	env.RegisterActivity(SendEmail)
	
	// Execute Workflow
	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)

}

func Test_FailedSubscriptionWorkflow (t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := EmailDetails{
		EmailAddress: "",
		Message: "This is a test to see if the Workflow will fail without a valid email.",
		IsSubscribed: false,
		SubscriptionCount: 0,
		MaxSubscriptionPeriods: 12,
	}

	env.RegisterWorkflow(SubscriptionWorkflow)
	env.RegisterActivity(SendEmail)
	
	// Execute Workflow
	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
}
// @@@SNIPEND