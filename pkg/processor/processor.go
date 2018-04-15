package processor

import (
	"reflect"
	"strings"
)

const nestingSeparator = ":"

const anyChildSelector = "*"
const anyArrayMemberSelector = "[]"

const actionAppend = "+"

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
			finalValue := append(typedTemplateValue, typedNewValue)
			data[head] = finalValue
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

func ProcessArraySelector(
	head interface{},
	remain []interface{},
	value interface{},
	data interface{},
) interface{} {
	remainLen := len(remain)

	switch typedHead := head.(type) {
	case string:
		switch typedHead {
		case anyArrayMemberSelector:
			switch typedLocation := data.(type) {
			case []interface{}:
				for index, member := range typedLocation {
					if remainLen == 0 {
						typedLocation[index] = value
					} else {
						nextHead := remain[0]
						nextRemain := remain[1:remainLen]

						switch typedMember := member.(type) {
						case map[interface{}]interface{}:
							typedLocation[index] = ProcessSelector(
								nextHead,
								nextRemain,
								value,
								typedMember,
							)
						}
					}
				}
			}
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
					nextHead := remain[0]
					nextRemain := remain[1:remainLen]

					switch nextData := data[key.Interface()].(type) {
					case map[interface{}]interface{}:
						data[key.Interface()] = ProcessSelector(
							nextHead,
							nextRemain,
							value,
							nextData,
						)
					}
				}
			}

		default:
			if remainLen == 0 {
				return doSet(head, value, data)
			} else if data[head] != nil {
				nextHead := remain[0]
				nextRemain := remain[1:remainLen]

				switch nextData := data[head].(type) {
				case map[interface{}]interface{}:
					data[head] = ProcessSelector(
						nextHead,
						nextRemain,
						value,
						nextData,
					)

				case []interface{}:
					data[head] = ProcessArraySelector(
						nextHead,
						nextRemain,
						value,
						nextData,
					)
				}
			}
		}
	}

	return data
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
			selectorParts := strings.Split(selectorStr, nestingSeparator)
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
