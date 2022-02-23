package tests

type A struct {
	A1 int
	A2 string
	A3 []int
	A4 [4]int
	A5 float64
	A6 map[string]interface{}
	A7 struct{}
}
type B struct {
	B1 int
	B2 string
	B3 [1]int
	B4 float64
}
type C struct {
	C1 []int
	C2 map[string]interface{}
	C3 struct{}
}

var D int

// ConvertToInterfaceSlice возвращает срез интерфейсов созданный из входных параметров
func ConvertToInterfaceSlice(args ...interface{}) []interface{} {
	var wants []interface{}
	for _, arg := range args {
		wants = append(wants, arg)
	}
	return wants
}
