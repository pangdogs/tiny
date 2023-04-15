package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func genEmit(ctx *_CommandContext) {
	emitFile := ctx.EmitDir

	if emitFile == "" {
		emitFile = strings.TrimSuffix(ctx.DeclFile, ".go") + "_emit_code.go"
	} else {
		emitFile = filepath.Dir(ctx.DeclFile) + string(filepath.Separator) + ctx.EmitDir + string(filepath.Separator) + filepath.Base(strings.TrimSuffix(ctx.DeclFile, ".go")) + "_emit_code.go"
	}

	emitCode := &bytes.Buffer{}

	// 生成注释
	{
		fmt.Fprintf(emitCode, `// Code generated by %s%s; DO NOT EDIT.

package %s
`, strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])),
			func() (args string) {
				for _, arg := range os.Args[1:] {
					if arg == "" {
						arg = `""`
					}
					args += " " + arg
				}
				return
			}(),
			ctx.EmitPackage)
	}

	// 生成import
	{
		importCode := &bytes.Buffer{}

		fmt.Fprintf(importCode, "\nimport (")

		fmt.Fprintf(importCode, `
	%s "kit.golaxy.org/tiny/localevent"`, ctx.EventPackageAlias)

		for _, imp := range ctx.FileAst.Imports {
			begin := ctx.FileSet.Position(imp.Pos())
			end := ctx.FileSet.Position(imp.End())

			impStr := string(ctx.FileData[begin.Offset:end.Offset])

			switch imp.Path.Value {
			case `"kit.golaxy.org/tiny/localevent"`:
				if imp.Name == nil {
					if ctx.EventPackageAlias == "localevent" {
						continue
					}
				} else {
					if imp.Name.Name == ctx.EventPackageAlias {
						continue
					}
				}
			case `"kit.golaxy.org/tiny/util"`:
				if imp.Name == nil {
					continue
				}
			}

			fmt.Fprintf(importCode, "\n\t%s", impStr)
		}

		fmt.Fprintf(importCode, `
	"kit.golaxy.org/tiny/util"`)

		fmt.Fprintf(importCode, "\n)\n")

		fmt.Fprintf(emitCode, importCode.String())
	}

	// 解析事件定义
	eventDeclTab := _EventDeclTab{}
	eventDeclTab.Parse(ctx)

	// localevent包前缀
	localeventPrefix := ""
	if ctx.EventPackageAlias != "." {
		localeventPrefix = ctx.EventPackageAlias + "."
	}

	// 生成事件发送代码
	for _, eventDecl := range eventDeclTab {
		// 是否导出事件发送代码
		exportEmitStr := "emit"
		if ctx.EmitDefExport {
			exportEmitStr = "Emit"
		}

		if strings.Contains(eventDecl.Comment, "[EmitExport]") {
			exportEmitStr = "Emit"
		} else if strings.Contains(eventDecl.Comment, "[EmitUnExport]") {
			exportEmitStr = "emit"
		}

		auto := ctx.EmitDefAuto

		if strings.Contains(eventDecl.Comment, "[EmitAuto]") {
			auto = true
		} else if strings.Contains(eventDecl.Comment, "[EmitManual]") {
			auto = false
		}

		// 生成代码
		if auto {
			if eventDecl.FuncHasRet {
				fmt.Fprintf(emitCode, `
type iAuto%[1]s interface {
	%[1]s() %[6]sIEvent
}

func Bind%[1]s(auto iAuto%[1]s, delegate %[2]s%[8]s) %[6]s.Hook {
	if auto == nil {
		panic("nil auto")
	}
	return %[6]sBindEvent[%[2]s%[8]s](auto.%[1]s(), delegate)
}

func Bind%[1]sWithPriority(auto iAuto%[1]s, delegate %[2]s%[8]s, priority int32) %[6]s.Hook {
	if auto == nil {
		panic("nil auto")
	}
	return %[6]sBindEventWithPriority[%[2]s%[8]s](auto.%[1]s(), delegate, priority)
}

func %[9]s%[1]s%[7]s(auto iAuto%[1]s%[4]s) {
	if auto == nil {
		panic("nil auto")
	}
	%[6]sUnsafeEvent(auto.%[1]s()).Emit(func(delegate util.IfaceCache) bool {
		return util.Cache2Iface[%[2]s%[8]s](delegate).%[3]s(%[5]s)
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, localeventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr)

			} else {
				fmt.Fprintf(emitCode, `
type iAuto%[1]s interface {
	%[1]s() %[6]sIEvent
}

func Bind%[1]s(auto iAuto%[1]s, delegate %[2]s%[8]s) %[6]s.Hook {
	if auto == nil {
		panic("nil auto")
	}
	return %[6]sBindEvent[%[2]s%[8]s](auto.%[1]s(), delegate)
}

func Bind%[1]sWithPriority(auto iAuto%[1]s, delegate %[2]s%[8]s, priority int32) %[6]s.Hook {
	if auto == nil {
		panic("nil auto")
	}
	return %[6]sBindEventWithPriority[%[2]s%[8]s](auto.%[1]s(), delegate, priority)
}

func %[9]s%[1]s%[7]s(auto iAuto%[1]s%[4]s) {
	if auto == nil {
		panic("nil auto")
	}
	%[6]sUnsafeEvent(auto.%[1]s()).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[%[2]s%[8]s](delegate).%[3]s(%[5]s)
		return true
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, localeventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr)
			}
		} else {
			if eventDecl.FuncHasRet {
				fmt.Fprintf(emitCode, `
func %[9]s%[1]s%[7]s(event %[6]sIEvent%[4]s) {
	if event == nil {
		panic("nil event")
	}
	%[6]sUnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		return util.Cache2Iface[%[2]s%[8]s](delegate).%[3]s(%[5]s)
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, localeventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr)

			} else {
				fmt.Fprintf(emitCode, `
func %[9]s%[1]s%[7]s(event %[6]sIEvent%[4]s) {
	if event == nil {
		panic("nil event")
	}
	%[6]sUnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[%[2]s%[8]s](delegate).%[3]s(%[5]s)
		return true
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, localeventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr)
			}
		}

		fmt.Printf("Emit: %s\n", eventDecl.Name)
	}

	os.MkdirAll(filepath.Dir(emitFile), os.ModePerm)

	if err := ioutil.WriteFile(emitFile, emitCode.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}
