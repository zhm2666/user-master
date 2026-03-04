package redis

import "strings"

const ServicePrefix = "user_crontab_"

func GetKey(key string, parts ...string) string {
	key = ServicePrefix + key
	if len(parts) == 0 {
		return key
	}
	key += "_" + strings.Join(parts, "_")
	return key
}
