-- KEYS:
-- [1]: token bucket key
-- [2]: last time key
-- ARGV:
-- [1]: interval_time
-- [2]: current_time
-- [3]: capacity
-- Return:
-- [1]: whether reach rate-limit（0: pass, 1: rate-limit）
-- [2]: current token amount

local res={}
res[1]=0

local interval_time=tonumber(ARGV[1]) -- time interval for putting token
local current_time=tonumber(ARGV[2]) -- current ts
local capacity=tonumber(ARGV[3])

local amount=1 -- token decrement each interval
local key_expire_time=1000*3600 -- EX
local inflow_per_unit=1 -- token increment each interval
local bucket_amount = 0

-- last time of putting token
local last_time=redis.call('get',KEYS[2])
-- current token amount
local current_value = redis.call('get',KEYS[1])

if(last_time == false or current_value == false) -- if token bucket not exist, create one
then
    bucket_amount = capacity - amount;
    -- gen token bucket
    redis.call('set',KEYS[1],bucket_amount,'EX',key_expire_time)
    -- set last time of putting token
    redis.call('set',KEYS[2],current_time,'EX',key_expire_time)
    res[2]=bucket_amount
    return res
end

current_value = tonumber(current_value)
last_time=tonumber(last_time)
local past_time=current_time-last_time -- time.since(last_time)
if(past_time<interval_time)
then
    -- if pass_time <  interval_time, no need to put token
    bucket_amount=current_value-amount
else
    local cur_times = past_time/interval_time
    cur_times=math.floor(cur_times)
    bucket_amount=current_value+cur_times*inflow_per_unit
    if (bucket_amount > capacity)
    then
        bucket_amount = capacity-amount
    end
    -- 有新投放，更新投放时间
    redis.call('set',KEYS[2],current_time,'PX',key_expire_time)
end

res[2]=bucket_amount

-- reach rate-limit
if(bucket_amount<0)
then
    res[1]=1
    return res
end

-- update token bucket KV
redis.call('set',KEYS[1],bucket_amount,'PX',key_expire_time)
return res
