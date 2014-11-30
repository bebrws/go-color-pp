package pp

import (
	"bytes"
	"fmt"
	. "github.com/k0kubun/palette"
	"reflect"
	"strings"
	"text/tabwriter"
)

const (
	indentWidth = 2
)

func format(object interface{}) string {
	return newPrinter(object).String()
}

func newPrinter(object interface{}) *printer {
	buffer := bytes.NewBufferString("")
	tw := new(tabwriter.Writer)
	tw.Init(buffer, indentWidth, 0, 1, ' ', 0)

	return &printer{
		Buffer: buffer,
		tw:     tw,
		depth:  0,
		value:  reflect.ValueOf(object),
	}
}

type printer struct {
	*bytes.Buffer
	tw    *tabwriter.Writer
	depth int
	value reflect.Value
}

func (p *printer) String() string {
	switch p.value.Kind() {
	case reflect.Bool:
		p.colorPrint(p.raw(), "Cyan")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Complex64, reflect.Complex128:
		p.colorPrint(p.raw(), "Blue")
	case reflect.Float32, reflect.Float64:
		p.colorPrint(p.raw(), "Magenta")
	case reflect.String:
		p.printString()
	case reflect.Map:
		p.printMap()
	case reflect.Struct:
		p.printStruct()
	default:
		p.print(p.raw())
	}

	p.tw.Flush()
	return p.Buffer.String()
}

func (p *printer) print(text string) {
	fmt.Fprint(p.tw, text)
}

func (p *printer) println(text string) {
	p.print(text + "\n")
}

func (p *printer) indentPrint(text string) {
	p.print(p.indent() + text)
}

func (p *printer) indentPrintf(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	p.indentPrint(text)
}

func (p *printer) colorPrint(text, color string) {
	p.print(Colorize(text, color))
}

func (p *printer) printString() {
	p.colorPrint(`"`, "Red")
	p.colorPrint(p.value.String(), "red")
	p.colorPrint(`"`, "Red")
}

func (p *printer) printMap() {
	p.println("{")
	p.indented(func() {
		keys := p.value.MapKeys()
		for i := 0; i < p.value.Len(); i++ {
			key := keys[i].Interface()
			value := p.value.MapIndex(keys[i]).Interface()
			p.indentPrintf("%s:\t%s,\n", p.format(key), p.format(value))
		}
	})
	p.indentPrint("}")
}

func (p *printer) printStruct() {
	t := p.value.Type()
	p.println(Green(t.String())+"{")
	p.indented(func() {
		for i := 0; i < p.value.NumField(); i++ {
			field := Yellow(t.Field(i).Name)
			value := p.value.Field(i).Interface()
			p.indentPrintf("%s:\t%s,\n", field, p.format(value))
		}
	})
	p.indentPrint("}")
}

func (p *printer) indented(proc func()) {
	p.depth++
	proc()
	p.depth--
}

func (p *printer) raw() string {
	return fmt.Sprintf("%#v", p.value.Interface())
}

func (p *printer) format(object interface{}) string {
	pp := newPrinter(object)
	pp.depth = p.depth
	return pp.String()
}

func (p *printer) indent() string {
	return strings.Repeat("\t", p.depth)
}