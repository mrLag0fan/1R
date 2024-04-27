package _1R

import (
	"1R/internal/scanner"
	"fmt"
	"reflect"
)

type Rule struct {
	AttrName string
	Attr     []interface{}
	Answer   []bool
	Accuracy float32
}

func Analyze(bestRule Rule, data *scanner.Data) {
	for i := range data.Data {
		entry := &data.Data[i]
		for j, attr := range bestRule.Attr {
			fieldValue := reflect.ValueOf(*entry).FieldByName(bestRule.AttrName)
			if fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.Bool {
				if fieldValue.Elem().Bool() == attr.(bool) {
					entry.Answer = &bestRule.Answer[j]
				}
			} else if fieldValue.IsValid() {
				if reflect.DeepEqual(fieldValue.Interface(), attr) {
					entry.Answer = &bestRule.Answer[j]
				}
			}
		}
	}
}

func FindBestRule(rules []Rule) *Rule {
	if len(rules) == 0 {
		return nil
	}

	bestRule := rules[0]
	for _, rule := range rules {
		if rule.Accuracy*100 > bestRule.Accuracy*100 {
			bestRule = rule
		}
	}

	return &bestRule
}

func Train1R(data scanner.Data) ([]Rule, error) {
	rules := make([]Rule, len(data.Attributes)-1)

	for i := 0; i < len(data.Attributes)-1; i++ {
		if i != 0 {
			rules[i] = *ruleForAttr(i, data)
		}
	}

	return rules, nil
}

func ruleForAttr(attr int, data scanner.Data) *Rule {
	fieldName := fieldNameByIndex(attr)
	attributes := setOfAttr(attr, data)
	fmt.Println(attributes)
	rule := &Rule{
		AttrName: fieldName,
		Attr:     attributes,
		Answer:   make([]bool, len(attributes)),
	}
	correct := 0
	if len(attributes) > 5 {
		// зробити розподіл на підмасиви і на основі них робити правила
		return nil
	} else {
		for i, attribute := range attributes {
			y, n := countAnswersForAttr(fieldName, attribute, data)
			fmt.Printf("For %v, Yes: %d, No: %d\n", attribute, y, n)
			if y > n {
				rule.Answer[i] = true
			} else {
				rule.Answer[i] = false
			}
			correct += y + n
		}
	}
	rule.Accuracy = calcAccuracy(fieldName, *rule, data)

	return rule
}

func calcAccuracy(attrName string, rule Rule, data scanner.Data) float32 {
	correct := 0
	for _, entry := range data.Data {
		for j, attr := range rule.Attr {
			fieldValue := reflect.ValueOf(entry).FieldByName(attrName)
			if fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.Bool {
				if fieldValue.Elem().Bool() == attr.(bool) && rule.Answer[j] == *entry.Answer {
					correct++
					break
				}
			}
			if reflect.DeepEqual(fieldValue.Interface(), attr) && rule.Answer[j] == *entry.Answer {
				correct++
				break
			}
		}
	}
	accuracy := float32(correct) / float32(len(data.Data))
	return accuracy
}

func countAnswersForAttr(attrName string, attr interface{}, data scanner.Data) (int, int) {
	countY := 0
	countN := 0
	for _, entry := range data.Data {
		fieldValue := reflect.ValueOf(entry).FieldByName(attrName)
		if fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.Bool {
			if fieldValue.Elem().Bool() == attr.(bool) {
				if *entry.Answer == true {
					countY++
				} else if *entry.Answer == false {
					countN++
				}
			}
		} else if fieldValue.IsValid() {
			if reflect.DeepEqual(fieldValue.Interface(), attr) {
				if *entry.Answer == true {
					countY++
				} else if *entry.Answer == false {
					countN++
				}
			}
		}
	}
	return countY, countN
}

func setOfAttr(attr int, data scanner.Data) []interface{} {
	res := make(map[interface{}]bool)
	for _, entry := range data.Data {
		addr := reflect.ValueOf(entry).FieldByName(fieldNameByIndex(attr)).Interface()
		key := addr
		if ptrBool, ok := addr.(*bool); ok {
			key = *ptrBool
		}
		res[key] = true
	}
	keys := make([]interface{}, 0, len(res))
	for key := range res {
		keys = append(keys, key)
	}

	return keys
}

func fieldNameByIndex(attr int) string {
	switch attr {
	case 0:
		return "Age"
	case 1:
		return "FrequencyPerYear"
	case 2:
		return "NewPlacesAttitude"
	case 3:
		return "HasPassport"
	}
	return ""
}
