package fixedlinemout_test

import (
	"testing"

	fixedlinemout "github.com/swaros/fixedline"
	"github.com/swaros/outinject"
)

func TestSplitByMarks(t *testing.T) {
	origin := "<width=55;lpad>part no 1<lpad>middle text<width=5;rpad>right text"
	parserEntrie, _ := fixedlinemout.NewRunner(false, "<>;=")

	parser := parserEntrie.ParseMarkup(origin)

	if len(parser.Entries) != 3 {
		t.Error("unexpected size of elements", len(parser.Entries), "  ", parser.Entries)
	}

	var testRunner fixedlinemout.Runner = fixedlinemout.Runner{Runners: make(map[string]fixedlinemout.MarkupRunner)}
	testRunner.HandleErrors = true

	if result, err := parser.ParseAll(testRunner); err != nil {
		t.Log(err, result)
	} else {
		t.Error("This have to be failing, because there are no runners defined ")
	}

	testRunner.Runners["width"] = fixedlinemout.MarkupRunner{
		Exec: func(mk fixedlinemout.Markup, current string) string {
			return "(:" + mk.Reference + ":)" + current
		},
	}

	testRunner.Runners["lpad"] = fixedlinemout.MarkupRunner{
		Exec: func(mk fixedlinemout.Markup, cur string) string {
			return "|:" + mk.Reference + ":|" + cur
		},
	}

	testRunner.Runners["rpad"] = fixedlinemout.MarkupRunner{
		Exec: func(mk fixedlinemout.Markup, cur string) string {
			return "):" + mk.Reference + ":(" + cur
		},
	}

	if res, err := parser.ParseAll(testRunner); err != nil {
		t.Error(err)
	} else {
		if res != "):right text:((:right text:)|:middle text:||:part no 1:|(:part no 1:)" {
			t.Error("wrong string:", res)
		}
	}

}

func benchmarkParsing(times int) {
	origin := "<width=55;lpad;lpad>check_1<lpad>check_1_1_1<width=5;rpad>S_Olskjdg_asgdfkgdfkaf_laksgfkhgf_lksajf_ohr"
	parserBase, _ := fixedlinemout.NewRunner(false, "<>;=")

	parser := parserBase.ParseMarkup(origin)

	parserBase.Runners["width"] = fixedlinemout.MarkupRunner{
		Exec: func(mk fixedlinemout.Markup, current string) string {
			return "(:" + mk.Reference + ":)" + current
		},
	}

	parserBase.Runners["lpad"] = fixedlinemout.MarkupRunner{
		Exec: func(mk fixedlinemout.Markup, cur string) string {
			return "|:" + mk.Reference + ":|" + cur
		},
	}

	parserBase.Runners["rpad"] = fixedlinemout.MarkupRunner{
		Exec: func(mk fixedlinemout.Markup, cur string) string {
			return "):" + mk.Reference + ":(" + cur
		},
	}
	for b := 0; b < times; b++ {
		parser.ParseAll(parserBase)
	}
}

func Benchmark(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkParsing(1)
	}
}

func TestSplitText(t *testing.T) {
	var setting fixedlinemout.FixedLine = fixedlinemout.FixedLine{
		LPanelMustShown: true,
		LPanelWidth:     10,
		RPanelWidth:     5,
		BufferWith:      0,
	}

	var fakeMo outinject.MOut = *outinject.NewStdout()
	fakeMo.IsTerminal = true
	fakeMo.Height = 10
	fakeMo.Width = 115

	enabled := setting.Enable(&fakeMo)
	if !enabled {
		t.Error("parser is not enabled")
	}

	leftString, middleString, rightString := setting.SplitSized("hello left panel", "hello message Content. i will see you hopefully", "XXXXXXX")

	if len(leftString) != setting.LPanelWidth {
		t.Error("string not matching expectations: '", leftString, "'", len(leftString))
	}

	if len(rightString) != setting.RPanelWidth {
		t.Error("string not matching expectations '", rightString, "'", len(rightString))
	}

	if len(middleString) != 100 {
		t.Error("string not matching expectations '", middleString, "'", len(middleString))
	}
}
