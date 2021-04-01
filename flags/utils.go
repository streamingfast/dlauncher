package flags

import (
	"fmt"
	"strings"
)

// ParseKVConfigString is used for flags that generate key/value data,
// 		                  like --backup="type=something freq_blocks=1000 prefix=v1"
func ParseKVConfigString(in string) (map[string]string, error) {
	fields := strings.Fields(in)
	kvs := map[string]string{}
	for _, field := range fields {
		kv := strings.Split(field, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid key=value in kv config string: %s", field)
		}
		kvs[kv[0]] = kv[1]
	}
	typ, ok := kvs["type"]
	if !ok || typ == "" {
		return nil, fmt.Errorf("no type defined in kv config string (type field mandatory)")
	}

	return kvs, nil
}
