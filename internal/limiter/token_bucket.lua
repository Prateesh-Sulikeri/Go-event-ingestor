-- Atomic Token Bucket with TTL
-- KEYS[1] bucket key
-- KEYS[2] timestamp key
-- ARGV[1] max_tokens
-- ARGV[2] refill_rate
-- ARGV[3] now

local bucket = KEYS[1]
local ts = KEYS[2]

local max = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local tokens = tonumber(redis.call("GET", bucket)) or max
local last = tonumber(redis.call("GET", ts)) or now

local elapsed = now - last
if elapsed < 0 then elapsed = 0 end
tokens = math.min(max, tokens + (elapsed * rate))

if tokens < 1 then
    redis.call("SETEX", bucket, 60, tokens)
    redis.call("SETEX", ts, 60, now)
    return 0
else
    tokens = tokens - 1
    redis.call("SETEX", bucket, 60, tokens)
    redis.call("SETEX", ts, 60, now)
    return 1
end
