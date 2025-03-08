package logic

import (
    "pyth-go/app/pyth-common/types"
    dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
)

func getAllDeduplicationKeys(service dtypes.IDeduplicationService, taskInfo *types.TaskInfo, tag string) (keys []string) {
	keys = make([]string, len(taskInfo.Receiver))

	for i, receiver := range taskInfo.Receiver {
		if tag != "" {
			keys[i] = tag + service.GetDeduplicationSingleKey(taskInfo, receiver)
		} else {
			keys[i] = service.GetDeduplicationSingleKey(taskInfo, receiver)
		}
	}
	return keys
}