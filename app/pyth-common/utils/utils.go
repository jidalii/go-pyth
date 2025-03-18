package utils

import (
	"fmt"

	"pyth-go/app/pyth-common/enums/chanType"
	"pyth-go/app/pyth-common/enums/msgType"
	"pyth-go/app/pyth-common/types"
)

const chanBusiness = "pythBusiness"

// Get groupId based on taskInfo
func GetGroupIdByTaskInfo(info types.TaskInfo) string {
	channelCodeEn := chanType.TypeCodeEn[info.SendChannel]
	msgCodeEn := msgType.TypeCodeEn[info.MsgType]
	return channelCodeEn + "." + msgCodeEn
}

// Get
func GetMqKey(channel, msgType string) string {
	return fmt.Sprintf("%s.%s.%s", chanBusiness, channel, msgType)
}

// Fetch all groupIds
// (different combinations of channel and msgType has a unique groupId)
func GetAllGroupIds() []string {
	chLen := len(chanType.TypeCodeEn)
	msgLen := len(msgType.TypeCodeEn)
	list := make([]string, chLen*msgLen)

	i := 0
	for _, ct := range chanType.TypeCodeEn {
		for _, mt := range msgType.TypeCodeEn {
			list[i] = ct + "." + mt
			i++
		}
	}
	return list
}
