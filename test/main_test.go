package code

import (
	"code"
	"testing"
	"github.com/stretchr/testify/require"
	format "code/formaters"
)
//json
func TestGenDiffEmptyfilesJson(t *testing.T) {
	m1:=code.Parsing("../testdata/json/file1_empty.json")
	m2:=code.Parsing("../testdata/json/file2_empty.json")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n}", ex1, "Empty files")
}

func TestGenDiffSimpleFilesWithSomeDifferencesJson(t *testing.T) {
	m1:=code.Parsing("../testdata/json/file1_simple.json")
	m2:=code.Parsing("../testdata/json/file2_simple.json")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n    a: 1\n  - b: test\n  + b: test2\n  - c: true\n  + d: false\n}", ex1, "FilesWithSomeDifferencesJson")
}
func TestGenDiffFirstFileEmptySecondWithDataJson(t *testing.T) {
	m1:=code.Parsing("../testdata/json/file1_empty.json")
	m2:=code.Parsing("../testdata/json/file1_simple.json")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n  + a: 1\n  + b: test\n  + c: true\n}", ex1, "FileEmptySecondWithDataJson")
}
func TestGenDiffFirstFileWithDataSecondEmptyJson(t *testing.T) {
	m1:=code.Parsing("../testdata/json/file1_simple.json")
	m2:=code.Parsing("../testdata/json/file1_empty.json")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n  - a: 1\n  - b: test\n  - c: true\n}", ex1, "FirstFileWithDataSecondEmptyJson")
}
func TestGenDiffmnJSON(t *testing.T) {
	m1:=code.Parsing("../testdata/json/file1_mn.json")
	m2:=code.Parsing("../testdata/json/file1_mn.json")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n    common: {\n    setting1: Value 1\n    setting2: 200\n    setting3: true\n    setting6: {\n    doge: {\n    wow: \n}\n    key: value\n}\n}\n    group1: {\n    baz: bas\n    foo: bar\n    nest: {\n    key: value\n}\n}\n    group2: {\n    abc: 12345\n    deep: {\n    id: 45\n}\n}\n}", ex1, "mnJSON")
}
//yml
func TestGenDiffEmptyfilesyml(t *testing.T) {
	m1:=code.Parsing("../testdata/yml/file1_empty.yml")
	m2:=code.Parsing("../testdata/yml/file2_empty.yml")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n}", ex1, "Empty files")
}

func TestGenDiffSimpleFilesWithSomeDifferencesyml(t *testing.T) {
	m1:=code.Parsing("../testdata/yml/file1_simple.yml")
	m2:=code.Parsing("../testdata/yml/file2_simple.yml")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n    a: 1\n  - b: test\n  + b: test2\n  - c: true\n  + d: false\n}", ex1, "FilesWithSomeDifferencesyml")
}
func TestGenDiffFirstFileEmptySecondWithDatayml(t *testing.T) {
	m1:=code.Parsing("../testdata/yml/file1_empty.yml")
	m2:=code.Parsing("../testdata/yml/file1_simple.yml")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n  + a: 1\n  + b: test\n  + c: true\n}", ex1, "FileEmptySecondWithDatayml")
}
func TestGenDiffFirstFileWithDataSecondEmptyyml(t *testing.T) {
	m1:=code.Parsing("../testdata/yml/file1_simple.yml")
	m2:=code.Parsing("../testdata/yml/file1_empty.yml")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n  - a: 1\n  - b: test\n  - c: true\n}", ex1, "FirstFileWithDataSecondEmptyyml")
}
func TestGenDiffmnyml(t *testing.T) {
	m1:=code.Parsing("../testdata/yml/file1_mn.yml")
	m2:=code.Parsing("../testdata/yml/file1_mn.yml")
    ex1:= code.GenDiff(m1,m2 )
    require.Equal(t, "{\n    common: {\n    setting1: Value 1\n    setting2: 200\n    setting3: true\n    setting6: {\n    doge: {\n    wow: \n}\n    key: value\n}\n}\n    group1: {\n    baz: bas\n    foo: bar\n    nest: {\n    key: value\n}\n}\n    group2: {\n    abc: 12345\n    deep: {\n    id: 45\n}\n}\n}", ex1, "mnyml")
}
//formats
//plan
func TestGenDiffSimplePlanJson(t *testing.T) {
	m1:=code.Parsing("../testdata/json/file1_simple.json")
	m2:=code.Parsing("../testdata/json/file2_simple.json")
    ex1:= format.FormatDiffOutput(code.GenDiff(m1,m2 ))
    require.Equal(t, "Property 'b' was added with value: 'test2'\nProperty 'c' was removed\nProperty 'd' was added with value: false", ex1, "TestGenDiffSimplePlanJson")
}
func TestGenDiffSimplePlanYml(t *testing.T) {
	m1:=code.Parsing("../testdata/yml/file1_simple.yml")
	m2:=code.Parsing("../testdata/yml/file2_simple.yml")
    ex1:= format.FormatDiffOutput(code.GenDiff(m1,m2 ))
    require.Equal(t, "Property 'b' was added with value: 'test2'\nProperty 'c' was removed\nProperty 'd' was added with value: false", ex1, "TestGenDiffSimplePlanYml")
}
//json
func TestGenDiffSimpleJsonJson(t *testing.T) {
	m1:=code.Parsing("../testdata/json/file1_simple.json")
	m2:=code.Parsing("../testdata/json/file2_simple.json")
    ex1:= format.FormatDiffToJSON(code.GenDiff(m1,m2 ))
    require.Equal(t, "{\n  \"common\": [\n    {\n      \"key\": \"b\",\n      \"value\": \"test2\"\n    },\n    {\n      \"key\": \"d\",\n      \"value\": false\n    }\n  ],\n  \"differences\": [\n    {\n      \"key\": \"b\",\n      \"newValue\": null,\n      \"oldValue\": \"test\",\n      \"type\": \"removed\"\n    },\n    {\n      \"key\": \"b\",\n      \"newValue\": \"test2\",\n      \"oldValue\": null,\n      \"type\": \"added\"\n    },\n    {\n      \"key\": \"c\",\n      \"newValue\": null,\n      \"oldValue\": true,\n      \"type\": \"removed\"\n    },\n    {\n      \"key\": \"d\",\n      \"newValue\": false,\n      \"oldValue\": null,\n      \"type\": \"added\"\n    }\n  ]\n}", ex1, "TestGenDiffSimpleJsonJson")
}
func TestGenDiffSimpleymlYml(t *testing.T) {
	m1:=code.Parsing("../testdata/yml/file1_simple.yml")
	m2:=code.Parsing("../testdata/yml/file2_simple.yml")
    ex1:= format.FormatDiffToJSON(code.GenDiff(m1,m2 ))
     require.Equal(t, "{\n  \"common\": [\n    {\n      \"key\": \"b\",\n      \"value\": \"test2\"\n    },\n    {\n      \"key\": \"d\",\n      \"value\": false\n    }\n  ],\n  \"differences\": [\n    {\n      \"key\": \"b\",\n      \"newValue\": null,\n      \"oldValue\": \"test\",\n      \"type\": \"removed\"\n    },\n    {\n      \"key\": \"b\",\n      \"newValue\": \"test2\",\n      \"oldValue\": null,\n      \"type\": \"added\"\n    },\n    {\n      \"key\": \"c\",\n      \"newValue\": null,\n      \"oldValue\": true,\n      \"type\": \"removed\"\n    },\n    {\n      \"key\": \"d\",\n      \"newValue\": false,\n      \"oldValue\": null,\n      \"type\": \"added\"\n    }\n  ]\n}", ex1, "TestGenDiffSimpleJsonYML")
}