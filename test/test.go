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

func ErrEqual(err1 error, err2 error) bool {
	return (err1 != nil && err2 != nil && err1.Error() == err2.Error()) ||
		(err1 == nil && err2 == nil)
}

func NewBool(value bool) *bool {
	b := value
	return &b
}
