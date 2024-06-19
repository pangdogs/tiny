package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func genEventTab(ctx *CommandContext) {
	code := &bytes.Buffer{}

	// 生成注释
	{
		program := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
		args := strings.Join(os.Args[1:], " ")

		fmt.Fprintf(code, `// Code generated by %s %s; DO NOT EDIT.

package %s
`, program, args, ctx.EventTabPackage)
	}

	// 生成import
	{
		importCode := &bytes.Buffer{}

		fmt.Fprintf(importCode, "\nimport (")

		fmt.Fprintf(importCode, `
	%s "%s"`, ctx.PackageEventAlias, packageEventPath)

		fmt.Fprintf(importCode, "\n)\n")

		fmt.Fprintf(code, importCode.String())
	}

	// 解析事件定义
	eventDeclTab := EventDeclTab{}
	eventDeclTab.Parse(ctx)

	// event包前缀
	eventPrefix := ""
	if ctx.PackageEventAlias != "." {
		eventPrefix = ctx.PackageEventAlias + "."
	}

	// 生成事件表接口
	{
		var eventsCode string

		for _, event := range eventDeclTab {
			eventsCode += fmt.Sprintf("\t%s() %sIEvent\n", event.Name, eventPrefix)
		}

		fmt.Fprintf(code, `
type I%[1]s interface {
%[2]s}
`, strings.Title(ctx.EventTabName), eventsCode)
	}

	// 生成事件表
	{
		var eventsRecursionCode string

		for i, event := range eventDeclTab {
			var eventRecursion string

			if strings.Contains(event.Comment, "[EventRecursion_Allow]") {
				eventRecursion = eventPrefix + "EventRecursion_Allow"
			} else if strings.Contains(event.Comment, "[EventRecursion_Disallow]") {
				eventRecursion = eventPrefix + "EventRecursion_Disallow"
			} else if strings.Contains(event.Comment, "[EventRecursion_Discard]") {
				eventRecursion = eventPrefix + "EventRecursion_Discard"
			} else if strings.Contains(event.Comment, "[EventRecursion_Truncate]") {
				eventRecursion = eventPrefix + "EventRecursion_Truncate"
			} else if strings.Contains(event.Comment, "[EventRecursion_Deepest]") {
				eventRecursion = eventPrefix + "EventRecursion_Deepest"
			} else {
				eventRecursion = "recursion"
			}

			eventsRecursionCode += fmt.Sprintf("\t(*eventTab)[%d].Init(autoRecover, reportError, %s, managed)\n", i, eventRecursion)
		}

		// 生成事件Id
		{
			fmt.Fprintln(code, `
var (`)
			fmt.Fprintf(code, `	_%[1]sId = %[2]sDeclareEventTabIdT[%[1]s]()
`, ctx.EventTabName, eventPrefix)

			for i, event := range eventDeclTab {
				fmt.Fprintf(code, `	%[2]sId = _%[1]sId + %[3]d
`, ctx.EventTabName, event.Name, i)
			}

			fmt.Fprintln(code, ")")
		}

		fmt.Fprintf(code, `
type %[1]s [%[2]d]%[4]sEvent

func (eventTab *%[1]s) Init(autoRecover bool, reportError chan error, recursion %[4]sEventRecursion, managed) {
%[3]s}

func (eventTab *%[1]s) Get(id uint64) %[4]sIEvent {
	if _%[1]sId != id & 0xFFFFFFFF00000000 {
		return nil
	}
	pos := id & 0xFFFFFFFF
	if pos >= uint64(len(*eventTab)) {
		return nil
	}
	return &(*eventTab)[pos]
}

func (eventTab *%[1]s) Open() {
	for i := range *eventTab {
		(*eventTab)[i].Open()
	}
}

func (eventTab *%[1]s) Close() {
	for i := range *eventTab {
		(*eventTab)[i].Close()
	}
}

func (eventTab *%[1]s) Clean() {
	for i := range *eventTab {
		(*eventTab)[i].Clean()
	}
}
`, ctx.EventTabName, len(eventDeclTab), eventsRecursionCode, eventPrefix)
	}

	for i, event := range eventDeclTab {
		fmt.Fprintf(code, `
func (eventTab *%[1]s) %[2]s() %[4]sIEvent {
	return &(*eventTab)[%[3]d]
}
`, ctx.EventTabName, event.Name, i, eventPrefix)
	}

	fmt.Printf("EventTab: %s\n", ctx.EventTabName)

	// 目标文件
	targetFile := filepath.Join(filepath.Dir(ctx.DeclFile), ctx.EventTabDir, filepath.Base(strings.TrimSuffix(ctx.DeclFile, ".go"))+"_tab_code.go")

	os.MkdirAll(filepath.Dir(targetFile), os.ModePerm)

	if err := ioutil.WriteFile(targetFile, code.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}
