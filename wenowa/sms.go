package wenowa

type SendSMSRequestStruct struct {
	Header      string `json:"header"`
	PhoneNumber string `json:"phoneNumber"`
	Message     string `json:"message"`
	Token       string `json:"token"`
	ScriptID    int64  `json:"scriptId,omitempty"`
	UsePackage  bool   `json:"usePackage"`
}

func (w Wenowa) SendSMS(req SendSMSRequestStruct) {
	
}
