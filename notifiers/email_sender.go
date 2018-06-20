package notifiers

import (
	"github.com/viktorminko/monitor/helper"
	"path"
)

func InitEmailSender(workDir string) (Sender, error) {
	dir := path.Join(
		path.Dir(workDir),
		path.Dir("notifiers/email/"),
	)

	gmailAccount := &EmailAccount{}
	err := helper.InitObjectFromJsonFile(path.Join(dir, "config.json"), gmailAccount)
	if err != nil {
		return nil, err
	}

	return &SmtpEmailSender{
		gmailAccount,
		dir,
	}, nil
}
