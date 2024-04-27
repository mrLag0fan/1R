package main

import (
	_1R "1R/internal/1R"
	"1R/internal/scanner"
	"fmt"
)

func main() {
	data := scanner.NewData("./data/data.csv")
	fmt.Println(data)
	rules, _ := _1R.Train1R(*data)
	for _, rule := range rules {
		fmt.Println(rule)
	}
	bestRule := _1R.FindBestRule(rules)
	fmt.Println(bestRule)
	toAnalyze := scanner.NewData("./data/analyze.csv")
	_1R.Analyze(*bestRule, toAnalyze)
	scanner.ExporyCSV(*toAnalyze, "./data/result.csv")
}
