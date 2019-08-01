package resources

import (
	"fmt"

	krsc "github.com/nimakaviani/kapp/pkg/kapp/resources"
)

type Conditions struct {
	Resource krsc.Resource
}

func (c Conditions) Reason(checkedType string) string {
	messages := c.reasons()

	if reason, found := messages[checkedType]; found {
		return reason
	}

	return "-"
}

func (c Conditions) IsSelectedTrue(checkedTypes []string) (bool, string) {
	statuses := c.statuses()

	for _, t := range checkedTypes {
		status, found := statuses[t]
		if !found {
			return false, ""
		}
		if status != "True" {
			return false, status
		}
	}

	return true, ""
}

func (c Conditions) IsAllTrue() (bool, string) {
	statuses := c.statuses()
	if len(statuses) == 0 {
		return false, "No conditions found"
	}

	for t, status := range c.statuses() {
		if status != "True" {
			return false, fmt.Sprintf("Condition %s is not True (%s)", t, status)
		}
	}

	return true, ""
}

func (c Conditions) statuses() map[string]string {
	statuses := map[string]string{}
	if conditions, ok := c.Resource.Status()["conditions"].([]interface{}); ok {
		for _, cond := range conditions {
			if typedCond, ok := cond.(map[string]interface{}); ok {
				if typedType, ok := typedCond["type"].(string); ok {
					if typedStatus, ok := typedCond["status"].(string); ok {
						statuses[typedType] = typedStatus
					}
				}
			}
		}
	}
	return statuses
}

func (c Conditions) reasons() map[string]string {
	statuses := map[string]string{}
	if conditions, ok := c.Resource.Status()["conditions"].([]interface{}); ok {
		for _, cond := range conditions {
			if typedCond, ok := cond.(map[string]interface{}); ok {
				if typedType, ok := typedCond["type"].(string); ok {
					if typedStatus, ok := typedCond["reason"].(string); ok {
						statuses[typedType] = typedStatus
					}
				}
			}
		}
	}
	return statuses
}
