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
	var activities *Activities

	testDetails := Subscription{
		EmailInfo{
			EmailAddress: "example@temporal.io",
			Mail: "",
		},
		Periods{
			TrialPeriod: time.Second,
			BillingPeriod: time.Second,
		},
	}

	env.RegisterWorkflow(SubscriptionWorkflow)

	env.RegisterActivity(activities.SendWelcomeEmail)
	env.RegisterActivity(activities.SendSubscriptionEmail)
	env.RegisterActivity(activities.SendCancellationEmail)
	env.RegisterActivity(activities.SendSubscriptionEndedEmail)
	
	// Execute Workflow
	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func Test_CanceledSubscriptionWorkflow (t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	var activities *Activities

	testDetails := Subscription{
		EmailInfo{
			EmailAddress: "example@temporal.io",
			Mail: "",
		},
		Periods{
			TrialPeriod: time.Second,
			BillingPeriod: time.Second,
		},
	}
	env.RegisterWorkflow(SubscriptionWorkflow)

	env.RegisterActivity(activities.SendWelcomeEmail)
	env.RegisterActivity(activities.SendSubscriptionEmail)
	env.RegisterActivity(activities.SendCancellationEmail)
	env.RegisterActivity(activities.SendSubscriptionEndedEmail)

	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
	env.CancelWorkflow()

}