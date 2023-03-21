package subscribe_emails

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Workflow definition
func SubscriptionWorkflow(ctx workflow.Context, subscription Subscription) error {

	var activities *Activities
	billingPeriodNum := 0
	duration := time.Minute

	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created for " + subscription.EmailInfo.EmailAddress)

	var queryResult string

	e := workflow.SetQueryHandler(ctx, "GetDetails", func(input []byte) (string, error) {
 		return queryResult, nil
	})
	if e != nil {
		logger.Info("SetQueryHandler failed: " + e.Error())
		return e 
	}

	var err error

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		WaitForCancellation: true,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// Handle any cleanup, including cancellations.
	defer func() {
		if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
			data := EmailInfo {
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail:         "Welcome! Looks like you've been signed up!",
			}
			e := workflow.ExecuteActivity(ctx, activities.SendCancellationEmailDuringActiveSubscription, data)

			if err != nil {
				logger.Error("Failed to send cancel email", "Error", e)
			} else {
				// Cancellation received, which will trigger an unsubscribe email.
				logger.Info("Sending cancellation email")
			}
			return
		}

		newCtx, _ := workflow.NewDisconnectedContext(ctx)

		data := EmailInfo {
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail: "You have been unsubscribed from the Subscription Workflow. Good bye.",
		}

		logger.Info("Sending unsubscribe email to " + subscription.EmailInfo.EmailAddress)
		err := workflow.ExecuteActivity(newCtx, activities.SendSubscriptionOverEmail, data).Get(newCtx, nil)

		if err != nil {
			logger.Error("Unable to send unsubscribe message", "Error", err)
		}
	}()

	logger.Info("Sending welcome email to " + subscription.EmailInfo.EmailAddress)

	data := EmailInfo {
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail:         "Welcome! Looks like you've been signed up!",
		}
			

	err = workflow.ExecuteActivity(ctx, activities.SendWelcomeEmail, data).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send welcome email", "Error", err)
	} else {
		billingPeriodNum++
	}

	// start subscription period
	for (billingPeriodNum < subscription.Periods.MaxBillingPeriods) {

		data := EmailInfo{
				EmailAddress: subscription.EmailInfo.EmailAddress,
				Mail:         "This is yet another email in the Subscription Workflow.",
		}

		err = workflow.ExecuteActivity(ctx, activities.SendSubscriptionEmail, data).Get(ctx, nil)

		if err != nil {
			logger.Error("Failed to send email ", "Error", err)
		}

		logger.Info("sent content email to " + subscription.EmailInfo.EmailAddress)

		err = workflow.ExecuteActivity(ctx, activities.ChargeCustomerForBillingPeriod, data).Get(ctx, nil)

		if err != nil {
			logger.Error("Failed to charge customer ", "Error", err)
		}

		billingPeriodNum++

		workflow.Sleep(ctx, duration)
	}

	return err
}

