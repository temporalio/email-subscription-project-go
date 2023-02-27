package subscribe_emails

import (
	"bufio"
	"context"
	"errors"
	"os"

	"go.temporal.io/sdk/activity"
)

// email activity
func SendContentEmail(ctx context.Context, emailInfo EmailInfo) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending email " + emailInfo.Mail + " to " + emailInfo.EmailAddress)

	message, err := getEmailFromFile(emailInfo.Mail)

	if err != nil {
		return sendMail(message, emailInfo.EmailAddress)
	}

	logger.Error("Failed getting email", err)
	return errors.New("unable to locate message to send")
}

func getEmailFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return scanner.Text(), scanner.Err()
}

func sendMail(message string, email string) error {
	return nil
}
