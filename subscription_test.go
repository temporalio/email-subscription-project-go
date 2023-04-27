package subscribe_emails

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_SuccessfulSubscriptionWorkflow (t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := Subscription{
		EmailInfo{
			EmailAddress: "example@temporal.io",
			Mail: "",
		},
		Periods{
			TrialPeriod: time.Second,
			BillingPeriod: time.Second,
			MaxBillingPeriods: 12,
		},
	}

	// Mock Activities
	
	// Execute Workflow
	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func Test_CanceledSubscriptionWorkflow (t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := Subscription{
		EmailInfo{
			EmailAddress: "example@temporal.io",
			Mail: "",
		},
		Periods{
			TrialPeriod: time.Second,
			BillingPeriod: time.Second,
			MaxBillingPeriods: 12,
		},
	}

	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}

func Test_FailedSubscriptionWorkflow (t *testing.T) {

}