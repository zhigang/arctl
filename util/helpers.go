package util

import (
	"fmt"
	"strings"
)

const (
	REQUEST_SCHEME   = "https"
	LABEL_PREFIX_JST = "label.jst.com/"
)

// ParseLabels returns a label maps and remove label arrays
func ParseLabels(labelSelector string) (map[string]string, []string, error) {
	spec := strings.Split(labelSelector, ",")
	labels := map[string]string{}
	var remove []string
	for _, labelSpec := range spec {
		if strings.Contains(labelSpec, "=") {
			parts := strings.Split(labelSpec, "=")
			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("invalid label spec: %v", labelSpec)
			}
			if strings.HasPrefix(parts[0], LABEL_PREFIX_JST) {
				labels[parts[0]] = parts[1]
			} else {
				labels[LABEL_PREFIX_JST+parts[0]] = parts[1]
			}
		} else if strings.HasSuffix(labelSpec, "-") {
			remove = append(remove, labelSpec[:len(labelSpec)-1])
		} else {
			return nil, nil, fmt.Errorf("unknown label spec: %v", labelSpec)
		}
	}
	for _, removeLabel := range remove {
		if _, found := labels[removeLabel]; found {
			return nil, nil, fmt.Errorf("can not both modify and remove a label in the same command")
		}
	}
	return labels, remove, nil
}

func GetBusinessMessage(businessCode string) string {
	switch businessCode {
	case "102000":
		return fmt.Sprintf("[%s]%s", businessCode, "没有发布准入凭据")
	case "102001":
		return fmt.Sprintf("[%s]%s", businessCode, "发布准入凭据审核中")
	default:
		return businessCode
	}
}
