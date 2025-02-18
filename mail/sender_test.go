package mail

import (
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
	`
	to := []string{"dinhtatuanlinh@gmail.com"}
	attachFiles := []string{"../test_send_file_through_email.txt"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
