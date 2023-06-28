package mail

import (
	"simple-bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Hello from Simple Bank"
	content := `
		<h1>Hello from Simple Bank</h1>
		<p>Your code is: <strong>123456</strong></p>
	`

	to := []string{"nguyenquocthang909@outlook.com.vn"}

	err = sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}