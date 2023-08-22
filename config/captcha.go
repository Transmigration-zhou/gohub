package config

import "gohub/pkg/config"

func init() {
	config.Add("captcha", func() map[string]interface{} {
		return map[string]interface{}{
			"height":            80,
			"width":             240,
			"length":            6,
			"max_skew":          0.7,
			"dot_count":         80,
			"expire_time":       15,
			"debug_expire_time": 10080,
			"testing_key":       "captcha_skip_test",
		}
	})
}
