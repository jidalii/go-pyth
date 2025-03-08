-- KEYS
-- [1]: sliding window rate-limiting key
-- ARGV
-- [1]: limit window (s)
-- [2]: current timestamp (served as score)
-- [3]: limit threshold
-- [4]: unique key for each score (placeholder)
-- Return
-- [0]: whether reach rate-limit (0: pass, 1: rate-limit)

local window = tonumber(ARGV[1])
local current_ts = tonumber(ARGV[2])
local threshold = tonumber(ARGV[3])

-- Remove the expired elements from the sorted set
redis.call('zremrangeByScore', KEYS[1], 0, current_ts - window)

-- calculate the current count
local current = redis.call('zcard', KEYS[1])

if (current == nil) or (current < threshold) then
    redis.call('zadd', KEYS[1], current_ts, ARGV[4])
    redis.call('expire', KEYS[1], window)
    return 0
else
    return 1
end