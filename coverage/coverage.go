package coverage

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
)

var out io.Writer
var seen map[string]bool

func init() {
	var err error
	out, err = os.OpenFile(os.Getenv("HUB_COVERAGE"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	seen = make(map[string]bool)
}

func Record(data interface{}, i int) {
	_, filename, _, _ := runtime.Caller(1)
	if !seen[filename] {
		seen[filename] = true
		d := reflect.ValueOf(data)
		count := reflect.ValueOf(d.FieldByName("Count").Interface())
		total := count.Len()
		for j := 0; j < total; j++ {
			write(data, j, 0, filename)
		}
	}
	write(data, i, 1, filename)
}

func write(data interface{}, i, count int, filename string) {
	d := reflect.ValueOf(data)
	pos := reflect.ValueOf(d.FieldByName("Pos").Interface())
	numStmt := reflect.ValueOf(d.FieldByName("NumStmt").Interface())

	fmt.Fprintf(
		out,
		"%s:%d.%d,%d.%d %d %d\n",
		filename,
		pos.Index(3*i).Uint(),
		pos.Index(3*i+2).Uint()&0xFFFF,
		pos.Index(3*i+1).Uint(),
		pos.Index(3*i+2).Uint()>>16&0xFFFF,
		numStmt.Index(i).Uint(),
		count,
	)
}
