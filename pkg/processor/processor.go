package processor

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Syntax constants
const nestingSeparator = ":"
const anyChildSelector = "*"
const anyArrayMemberSelector = "[]"
const actionAppend = "+"

// Syntax regexes
var arrayMemberSelector = regexp.MustCompile("^\\[([0-9]*)\\](.)?$")
var arraySearcherSelector = regexp.MustCompile("^\\[([^=]+)=([^\\]]+)\\](.)?$")

func ProcessArraySelector(
	head interface{},
	remain []interface{},
	value interface{},
	data interface{},
) interface{} {
	remainLen := len(remain)

	switch typedHead := head.(type) {
	case string:
		switch typedLocation := data.(type) {
		case []interface{}:

			// This has to be a prefix check because action characters could follow
			// the any array selector if its in the sentinal position.
			if strings.HasPrefix(typedHead, anyArrayMemberSelector) {
				for index := range typedLocation {
					if remainLen == 0 {
						data = doSetArray(head, value, index, typedLocation)
					} else {
						nextHead := remain[0]
						nextRemain := remain[1:remainLen]

						switch typedMember := typedLocation[index].(type) {
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
			} else if arrayMemberSelector.MatchString(typedHead) {
				arrayMemberSelectorMatches := arrayMemberSelector.FindStringSubmatch(typedHead)
				selectedIndex, _ := strconv.Atoi(arrayMemberSelectorMatches[1])

				if remainLen == 0 {
					data = doSetArray(head, value, selectedIndex, typedLocation)
				} else {
					nextHead := remain[0]
					nextRemain := remain[1:remainLen]

					switch typedMember := typedLocation[selectedIndex].(type) {
					case map[interface{}]interface{}:
						typedLocation[selectedIndex] = ProcessSelector(
							nextHead,
							nextRemain,
							value,
							typedMember,
						)
					}
				}
			} else if arraySearcherSelector.MatchString(typedHead) {
				arraySearcherSelectorMatches := arraySearcherSelector.FindStringSubmatch(typedHead)
				searchProperty := interface{}(arraySearcherSelectorMatches[1])
				searchValue := interface{}(arraySearcherSelectorMatches[2])
				//actionCharacter := arraySearcherSelectorMatches[3]
				nextHead := remain[0]
				nextRemain := remain[1:remainLen]

				for index := range typedLocation {
					switch typedMember := typedLocation[index].(type) {
					case map[interface{}]interface{}:
						if typedMember[searchProperty] == searchValue {
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
