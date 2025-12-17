package redisutil

const (
	allocateScript = `
if redis.call("EXISTS", KEYS[1]) == 0 then
	redis.call("SET", KEYS[1], ARGV[1])
	redis.call("EXPIRE", KEYS[1], ARGV[2])
	return 1
else
	return 0
end
`
	renewScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
redis.call("EXPIRE", KEYS[1], ARGV[2])
return 1
else
return 0
end
`
)
