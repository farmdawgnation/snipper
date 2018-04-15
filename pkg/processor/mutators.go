package processor

import (
  "strings"
)

/**
 * An append that won't tinker with its input parameters.
 */
func safeAppend(source []interface{}, additional ...interface{}) (result []interface{}) {
  totalLen := len(source) + len(additional)

  result = make([]interface{}, totalLen)
  result = append(result, source)
  result = append(result, additional)
  return
}

func doAppend(head interface{}, value interface{}, data map[interface{}]interface{}) map[interface{}]interface{} {
	switch typedTemplateValue := data[head].(type) {
	case string:
		switch typedNewValue := value.(type) {
		case string:
			data[head] = typedTemplateValue + typedNewValue
		}

	case []interface{}:
		switch typedNewValue := value.(type) {
		case map[interface{}]interface{}:
			finalValue := safeAppend(typedTemplateValue, typedNewValue)
			data[head] = finalValue
		case []interface{}:
			finalValue := safeAppend(typedTemplateValue, typedNewValue...)
			data[head] = finalValue
		}
	}


	return data
}

func doSet(head interface{}, value interface{}, data map[interface{}]interface{}) map[interface{}]interface{} {
	switch typedHead := head.(type) {
	case string:
		if strings.HasSuffix(typedHead, actionAppend) {
			typedHeadWithoutSuffix := strings.TrimSuffix(typedHead, actionAppend)
			var newHeadIntf interface{}
			newHeadIntf = typedHeadWithoutSuffix

			data = doAppend(newHeadIntf, value, data)
		} else {
			data[head] = value
		}
	}

	return data
}

func doAppendArray(head interface{}, value interface{}, index int, data []interface{}) []interface{} {
	switch typedTemplateValue := data[index].(type) {
	case string:
		switch typedNewValue := value.(type) {
		case string:
			data[index] = typedTemplateValue + typedNewValue
		}
	}


	return data
}

func doSetArray(head interface{}, value interface{}, index int, data []interface{}) []interface{} {
	switch typedHead := head.(type) {
	case string:
		if strings.HasSuffix(typedHead, actionAppend) {
			typedHeadWithoutSuffix := strings.TrimSuffix(typedHead, actionAppend)
			var newHeadIntf interface{}
			newHeadIntf = typedHeadWithoutSuffix

			doAppendArray(newHeadIntf, value, index, data)
		} else {
			data[index] = value
		}
	}

	return data
}
