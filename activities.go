// @@@SNIPSTART subscription-workflow-go-activities
package subscribe_emails

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
)

// email activities
func SendEmail(ctx context.Context, emailInfo EmailDetails) (string, error) {
	activity.GetLogger(ctx).Info("Sending email to the following customer: %v", emailInfo.EmailAddress)
	return fmt.Sprintf("Sent email %v out of %v to: %v", emailInfo.SubscriptionCount, emailInfo.MaxSubscriptionPeriods, emailInfo.EmailAddress), nil
}
// @@@SNIPEND
