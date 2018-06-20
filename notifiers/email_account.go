package notifiers

type Recipient struct {
	Name string
	Email string
}

type EmailAccount struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
	IsDebugMode bool
	Recipients []Recipient
	From string
}

func (a *EmailAccount) GetRecipients() ([]string, error) {
	var emails []string
	for _, recipient := range a.Recipients {
		emails = append(emails, recipient.Email)
	}

	return emails, nil
}
