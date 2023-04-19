package subscribe_emails

import (
	"testing"
	"time"

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
			BillingPeriodCharge: 10,
		},
	}

	// Mock Activities
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
			BillingPeriodCharge: 10,
		},
	}
}

func Test_FailedSubscriptionWorkflow (t *testing.T) {

}