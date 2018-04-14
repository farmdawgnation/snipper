package processor

import (
	"reflect"
	"strings"
)

const nestingSeparator = "."

const anyChildSelector = "*"

const actionAppend = "+"

func doAppend(head interface{}, value interface{}, data map[interface{}]interface{}) map[interface{}]interface{} {
	switch typedTemplateValue := data[head].(type) {
	case string:
		switch typedNewValue := value.(type) {
		case string:
			finalValue := typedTemplateValue + typedNewValue
			data[head] = finalValue
		}

	case []interface{}:
		switch typedNewValue := value.(type) {
		case []interface{}:
			finalValue := append(typedTemplateValue, typedNewValue...)
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

/**
 * Process an individual selector value
 */
func ProcessSelector(
	head interface{},
	remain []interface{},
	value interface{},
	data map[interface{}]interface{},
) map[interface{}]interface{} {
	remainLen := len(remain)

	switch typedHead := head.(type) {
	case string:
		switch typedHead {
		case anyChildSelector:
			keys := reflect.ValueOf(data).MapKeys()
			for _, key := range keys {
				if remainLen == 0 {
					data = doSet(key.Interface(), value, data)
				} else {
					switch nextData := data[key.Interface()].(type) {
					case map[interface{}]interface{}:
						nextHead := remain[0]
						nextRemain := remain[1:remainLen]
						data[key.Interface()] = ProcessSelector(
							nextHead,
							nextRemain,
							value,
							nextData,
						)
					}
				}
			}

			return data

		default:
			if remainLen == 0 {
				return doSet(head, value, data)
			} else if data[head] != nil {
				switch nextData := data[head].(type) {
				case map[interface{}]interface{}:
					nextHead := remain[0]
					nextRemain := remain[1:remainLen]
					data[head] = ProcessSelector(
						nextHead,
						nextRemain,
						value,
						nextData,
					)
					return data
				}
			}
		}
	}

	return nil
}

/**
 * Run the provided template through the provided transform.
 */
func Process(
	template map[interface{}]interface{},
	transform map[interface{}]interface{},
) map[interface{}]interface{} {
	selectors := reflect.ValueOf(transform).MapKeys()

	for _, selector := range selectors {
		selectorIntf := selector.Interface()
		selectorSetting := transform[selectorIntf]

		switch selectorStr := selectorIntf.(type) {
		case string:
			selectorParts := strings.Split(selectorStr, ".")
			numParts := len(selectorParts)
			selectorHead := selectorParts[0]
			var selectorHeadIntf interface{}
			selectorHeadIntf = selectorHead
			selectorRemain := selectorParts[1:numParts]

			selectorRemainIntf := make([]interface{}, len(selectorRemain))
			for index, value := range selectorRemain {
				selectorRemainIntf[index] = value
			}

			template = ProcessSelector(selectorHeadIntf, selectorRemainIntf, selectorSetting, template)
		}


	}

	return template
}
