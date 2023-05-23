// @@@SNIPSTART subscription-workflow-go-subscribe-test
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
		EmailInfo {
			EmailAddress: "example@temporal.io",
			Mail: "",
		},
		SubscriptionPeriod: ,
		MaxSubscriptionPeriods: 12,
	}

	env.RegisterWorkflow(SubscriptionWorkflow)

	env.RegisterActivity(activities.SendWelcomeEmail)
	env.RegisterActivity(activities.SendSubscriptionEmail)
	env.RegisterActivity(activities.SendCancellationEmail)
	env.RegisterActivity(activities.SendSubscriptionEndedEmail)
	
	// Execute Workflow
	env.ExecuteWorkflow(SubscriptionWorkflow, testDetails)
	require.NoError(t, env.GetWorkflowError())
}

// @@@SNIPEND