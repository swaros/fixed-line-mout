package fixedlinemout_test

import (
	"testing"

	fixedlinemout "github.com/swaros/fixedline"
	"github.com/swaros/outinject"
)

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
