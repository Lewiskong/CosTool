package tfmt

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
)

func (tf *TerminalFormatter) Output() {
	var result strings.Builder

	// output header
	for i, field := range tf.fields {
		result.WriteString(field + spaces(tf.maxFieldLen[i]-len(field)) + "\t")
	}
	result.WriteString("\n")

	// output body
	for _, line := range tf.lines {
		for i, fieldInterface := range line {
			field := cast.ToString(fieldInterface)
			result.WriteString(field + spaces(tf.maxFieldLen[i]-len(field)) + "\t")
		}
		result.WriteString("\n")
	}
	fmt.Print(result.String())
}

func spaces(num int) string {
	sb := strings.Builder{}
	for i := 0; i < num; i++ {
		sb.WriteString(" ")
	}
	return sb.String()
}

func (tf *TerminalFormatter) Println(args ...interface{}) {
	if len(args) != len(tf.fields) {
		return
	}
	tf.lines = append(tf.lines, args)
	for i, arg := range args {
		s := cast.ToString(arg)
		if tf.maxFieldLen[i] < len(s) {
			tf.maxFieldLen[i] = len(s)
		}
	}
}

func New(fields ...string) *TerminalFormatter {
	tf := &TerminalFormatter{}
	for _, f := range fields {
		tf.fields = append(tf.fields, f)
		tf.maxFieldLen = append(tf.maxFieldLen, len(f))
	}
	return tf
}

type TerminalFormatter struct {
	fields      []string
	lines       [][]interface{}
	maxFieldLen []int
}
