package engine

import (
	"fmt"
	"os"
	"testing"

	"github.com/dop251/goja"
)

func TestGojaAdd(t *testing.T) {
	filePath := "test_data/add.js" // os.Args[1]
	jsData, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	vm := goja.New()
	v, err := vm.RunString(string(jsData))
	if err != nil {
		panic(err)
	}

	// Export：将结果转换成go基础类型
	result, ok := v.Export().(int64)
	fmt.Printf("result: %d, ok: %t\n", result, ok)
}

func TestGojaMapStruct(t *testing.T) {

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	type S struct {
		Field int `json:"field"`
	}
	vm.Set("s", S{Field: 42})
	res, _ := vm.RunString(`s.field`) // without the mapper it would have been s.Field
	fmt.Println(res.Export())
	// Output: 42
}
