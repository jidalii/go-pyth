package types

type TaskInfo struct {
	MessageTemplateId int64       `json:"messageTemplateId"`
	BusinessId        int64       `json:"businessId"`
	Receiver          []string    `json:"receiver"`
	IdType            int         `json:"idType"`
	SendChannel       int         `json:"sendChannel"`
	TemplateType      int         `json:"templateType"`
	MsgType           int         `json:"msgType"`
	ShieldType        int         `json:"shieldType"`
	ContentModel      interface{} `json:"contentModel"`
	SendAccount       int64       `json:"sendAccount"`
}

type SendTaskModel struct {
	MessageTemplateId int64          `json:"messageTemplateId"`
	MessageParamList  []MessageParam `json:"messageParamList"`
	TaskInfo          []TaskInfo     `json:"taskInfo"`
}

type MessageParam struct {
	Receiver  string                 `json:"receiver"`           // Receivers if multiple, separated by commas
	Variables map[string]interface{} `json:"variables,optional"` // Optional 消息内容中的可变部分(占位符替换)
	Extra     map[string]interface{} `json:"extra,optional"`     // Optional 扩展参数
}
