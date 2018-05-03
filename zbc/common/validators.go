package zbcommon

import (
	"strconv"
	"strings"
)

func IsIPv4(hostport string) bool {
	hostPort := strings.Split(hostport, ":")
	if len(hostPort) != 2 {
		return false
	}
	host := hostPort[0]
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return false
	}
	if port < 1001 || port > 65536 {
		return false
	}

	parts := strings.Split(host, ".")
	if len(parts) < 4 {
		return false
	}

	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}

	}
	return true
}
