package test

type OrderedTests struct {
	testDataSet TestDataSet
	orderedList OrderedTestList
}

type TestDataSet map[string]TestData
type OrderedTestList []string

type TestData struct {
	expected interface{}
	data     interface{}
}

func ErrEqual(err1 error, err2 error) bool {
	return (err1 != nil && err2 != nil && err1.Error() != err2.Error()) ||
		(err1 == nil && err1 != err2) ||
		(err2 == nil && err1 != err2)
}
