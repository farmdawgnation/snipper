package processor

import "testing"
import "github.com/farmdawgnation/snipper/pkg/processor"

func TestSimpleValueSet(t *testing.T) {
  dataMap := make(map[interface{}]interface{})
  remainList := make([]interface{}, 0)

  dataMap["bacon"] = "bacon"

  resultingMap := processor.ProcessSelector("bacon", remainList, "walnut", dataMap)

  if resultingMap["bacon"] != "walnut" {
    t.Error("Expected walnut got", resultingMap["bacon"])
  }
}

func TestNestedValueSet(t *testing.T) {
  dataMap := make(map[interface{}]interface{})
  innerDataMap := make(map[interface{}]interface{})

  remainList := make([]interface{}, 1)
  remainList[0] = "bacon"

  innerDataMap["bacon"] = "bacon"
  dataMap["parent"] = innerDataMap

  resultingMap := processor.ProcessSelector("parent", remainList, "walnut", dataMap)

  switch innerMap := resultingMap["parent"].(type) {
  case map[interface{}]interface{}:
    if innerMap["bacon"] != "walnut" {
      t.Error("Expected walnut got ", innerMap["bacon"])
    }

  default:
    t.Error("Unexpected type under parent key")
  }
}

func TestWildcardValueSet(t *testing.T) {
  dataMap := make(map[interface{}]interface{})
  middleDataMap := make(map[interface{}]interface{})
  innerDataMap1 := make(map[interface{}]interface{})
  innerDataMap2 := make(map[interface{}]interface{})

  innerDataMap1["name"] = "Shadow"
  innerDataMap2["name"] = "Shadow"

  middleDataMap["dog1"] = innerDataMap1
  middleDataMap["dog2"] = innerDataMap2

  dataMap["dogs"] = middleDataMap

  remainList := make([]interface{}, 2)
  remainList[0] = "*"
  remainList[1] = "name"

  resultingMap := processor.ProcessSelector("dogs", remainList, "Beamer", dataMap)

  switch middleMap := resultingMap["dogs"].(type) {
  case map[interface{}]interface{}:
    switch rInnerMap1 := middleMap["dog1"].(type) {
    case map[interface{}]interface{}:
      if rInnerMap1["name"] != "Beamer" {
        t.Error("Was expecting Beamer in map 1 got", rInnerMap1["name"])
      }

    default:
      t.Error("Got unexpected type on inner map 1")
    }

    switch rInnerMap2 := middleMap["dog2"].(type) {
    case map[interface{}]interface{}:
      if rInnerMap2["name"] != "Beamer" {
        t.Error("Was expecting Beamer in map 2 got", rInnerMap2["name"])
      }

    default:
      t.Error("Got unexpected type on inner map 2")
    }

  default:
    t.Error("Unexpected type under dogs key")
  }
}