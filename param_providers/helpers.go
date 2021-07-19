package param_providers

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func FilterParams(anyParams map[string]interface{}, wantedKeys []string) map[string]interface{} {
	var ret = make(map[string]interface{})
	for k, v := range anyParams {
		if contains(wantedKeys, k) {
			ret[k] = v
		}
	}
	return ret
}

func ForceMapValuesToString(anyParams map[string]interface{}) map[string]string {
	var ret = make(map[string]string)
	for k, v := range anyParams {
		ret[k] = fmt.Sprintf("%v", v)
	}
	return ret
}

func HashParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// collect map keys & values into one long string, in k,v,k,v,... order
	blob := ""

	for _, k := range keys {
		blob += fmt.Sprintf("%s%s", k, params[k])
	}

	sum := sha256.Sum256([]byte(blob))
	return fmt.Sprintf("%x", sum)
}
