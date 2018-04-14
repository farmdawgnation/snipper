package processor

import (
	//"fmt"
	"reflect"
	"strings"
)

const nestingSeparator = "."

const anyChildSelector = "*"

const actionSet = "="
const actionAppend = "+"

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
					data[key.Interface()] = value
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
				data[head] = value
				return data
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
