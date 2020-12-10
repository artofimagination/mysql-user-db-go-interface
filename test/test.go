package test

type OrderedTests struct {
	TestDataSet DataSet
	OrderedList OrderedTestList
}

type DataSet map[string]Data
type OrderedTestList []string

type Data struct {
	Expected interface{}
	Data     interface{}
}

var TestResultString = "\n%s test failed.\n\nReturned:\n%+v\n\nExpected:\n%+v"
