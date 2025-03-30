package utils

import (
	"fmt"
	"strings"

	"pyth-go/app/pyth-common/enums/chanType"
	"pyth-go/app/pyth-common/enums/msgType"
	"pyth-go/app/pyth-common/types"
)

const chanBusiness = "pythBusiness"

//	ReplaceByMap returns a copy of `origin`,
//
// which is replaced by a map in unordered way, case-sensitively.
func ReplaceByMap(origin string, replaces map[string]string) string {
	for k, v := range replaces {
		origin = strings.Replace(origin, "{$"+k+"}", v, -1)
	}
	return origin
}

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
