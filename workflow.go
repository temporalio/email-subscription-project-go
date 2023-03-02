package subscribe_emails

import (
	"go.temporal.io/sdk/workflow"
)

// Workflow definition
func SubscriptionWorkflow(ctx workflow.Context, subscription Subscription) error {
	// declare variables
	workflowCustomer := subscription
	subscriptionCancelled := false
	billingPeriodNum := 0
	actResult := ""

	QueryCustomerIdName := "customerid"
	QueryBillingPeriodNumberName := "billingperiodnumber"
	QueryBillingPeriodChargeAmountName := "billingperiodchargeamount"

	// set up logger
	logger := workflow.GetLogger(ctx)

	// query handler
	err := workflow.SetQueryHandler(ctx, QueryCustomerIdName, func() (string, error) {
		return workflowCustomer.EmailInfo.FirstName, nil
	})
	if err != nil {
		logger.Info("QueryCustomerIdName handler failed.", "Error", err)
		return "Error", err
	}

	err = workflow.SetQueryHandler(ctx, QueryBillingPeriodNumberName, func() (int, error) {
		return billingPeriodNum, nil
	})
	if err != nil {
		logger.Info("QueryBillingPeriodNumberName handler failed.", "Error", err)
		return "Error", err
	}

	err = workflow.SetQueryHandler(ctx, QueryBillingPeriodChargeAmountName, func() (int, error) {
		return workflowCustomer.Subscription.BillingPeriodCharge, nil
	})
	if err != nil {
		logger.Info("QueryBillingPeriodChargeAmountName handler failed.", "Error", err)
		return "Error", err
	}

	// handle subscription with handlers

	// send welcome email, start free trial

	// handle cancellation

	// handle unsubscription

	return "Completed Subscription Workflow", err
}
