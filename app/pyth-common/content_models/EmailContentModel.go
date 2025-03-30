package content_model

import (
	"github.com/zeromicro/go-zero/core/jsonx"
	"pyth-go/app/pyth-common/model"
	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-common/utils"
)

type EmailContentModel struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewEmailContentModel() *EmailContentModel {
	return &EmailContentModel{}
}

func (d EmailContentModel) BuilderContent(messageTemplate model.MessageTemplate, messageParam types.MessageParam) interface{} {
	variables := messageParam.Variables
	var content EmailContentModel
	_ = jsonx.Unmarshal([]byte(messageTemplate.MsgContent), &content)
	newVariables := getStringVariables(variables)
	content.Content = utils.ReplaceByMap(content.Content, newVariables)
	content.Title = newVariables["title"]
	return content
}

func getStringVariables(variables map[string]interface{}) map[string]string {
	var newVariables = make(map[string]string, len(variables))
	for key, variable := range variables {
		if v, ok := variable.(string); ok {
			newVariables[key] = v
		}
	}
	return newVariables
}
