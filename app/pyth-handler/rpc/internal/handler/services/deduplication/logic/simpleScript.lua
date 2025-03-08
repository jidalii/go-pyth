-- Params:
-- KEYS: ["SP_key1", "SP_key2", ...], where "SP_" is simpleLimitServiceTag, and key1、key2 is the result of deduplicationSingleKey
-- ARGV[1]: threshold（param.Num）
-- ARGV[2]: expiration (param.Time) in seconds
-- ARGV[3...]: receiver list, the order aligns with KEYS

local threshold = tonumber(ARGV[1])
local expiration = tonumber(ARGV[2])
local filteredReceiver = {}
local n = #KEYS

-- execute MGET to get all keys' values
local values = redis.call("MGET", unpack(KEYS))

for i = 1, n do
    local receiver = ARGV[i + 2]  -- ARGV[3] 对应 KEYS[1]
    local value = values[i]
    if not value then
        -- If key not exist, init to 1 and set EX
        redis.call("SET", KEYS[i], 1, "EX", expiration)
    else
        local count = tonumber(value)
        if count > threshold then
            -- If already reach threshold, add receiver to filtered list
            table.insert(filteredReceiver, receiver)
        else
            -- Otherwise, update count and EX
            count = count + 1
            redis.call("SET", KEYS[i], count, "EX", expiration)
        end
    end
end

-- return filtered receiver list
return filteredReceiver
