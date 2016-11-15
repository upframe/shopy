package main

type config struct {
	Development    bool   `json:"Development"`
	Key            string `json:"Key"`
	InviteOnly     bool   `json:"InviteOnly"`
	DefaultInvites int    `json:"BaseInvites"`
	Database       struct {
		User     string `json:"User"`
		Password string `json:"Password"`
		Host     string `json:"Host"`
		Port     string `json:"Port"`
		Name     string `json:"Name"`
	} `json:"Database"`
	SMTP struct {
		User     string `json:"User"`
		Password string `json:"Password"`
		Host     string `json:"Host"`
		Port     string `json:"Port"`
	} `json:"SMTP"`
	PayPal struct {
		Client string `json:"Client"`
		Secret string `json:"Secret"`
	} `json:"PayPal"`
}
