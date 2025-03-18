package chanType

const (
	Im    int = 10 // Internal Message
	Push  int = 20 // Push Notification
	Sms   int = 30 // SMS
	Email int = 40 // Email
)

type ChanMetadata struct {
	Description       string
	CodeEn            string
	AccessTokenPrefix string
	AccessTokenExpire int
}

var (
	Metadata = map[int]ChanMetadata{
		Im:    {"Internal Message", "im", "", -1},
		Push:    {"Push Notification", "push", "ge_tui_access_token_", 3600*24},
		Sms:    {"SMS", "sms", "", -1},
		Email:    {"Email", "email", "", -1},
	}
)

var (
	TypeText = map[int]string{
		Im:    "IM",
		Push:  "PUSH",
		Sms:   "SMS",
		Email: "EMAIL",
	}
	TypeCodeEn = map[int]string{
		Im:    "im",
		Push:  "push",
		Sms:   "sms",
		Email: "email",
	}
)
