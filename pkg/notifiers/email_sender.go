package notifiers

import (
	"github.com/viktorminko/monitor/pkg/helper"
	"path"
)

func InitEmailSender(workDir string) (Sender, error) {
	dir := path.Join(workDir, "notifiers/email")

	gmailAccount := &EmailAccount{}
	err := helper.InitObjectFromJsonFile(dir, "config.json", gmailAccount)
	if err != nil {
		return nil, err
	}

	return &SmtpEmailSender{
		gmailAccount,
		dir,
	}, nil
}
