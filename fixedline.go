/*MIT License

Copyright (c) 2022 Thomas Ziegler, <thomas.zglr@googlemail.com>. All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package fixedlinemout

import (
	"fmt"

	"github.com/swaros/outinject"
)

type FixedLine struct {
	enabled         bool
	width           int
	height          int
	haveSize        bool
	LPanelMustShown bool
	LPanelWidth     int
	RPanelWidth     int
	BufferWith      int
}

const (
	M_FL_LEFT_ALIGN  = "<fl:left>"
	M_FL_RIGHT_ALIGN = "<fl:right>"
)

func (f *FixedLine) Enable(mo *outinject.MOut) bool {
	f.enabled = mo.IsTerminal
	f.height = 0
	f.width = 0
	f.haveSize = false
	if f.enabled && mo.Width > 0 {
		f.height = mo.Height
		f.width = mo.Width
		f.haveSize = true
	}
	return f.enabled

}

func (f *FixedLine) Parse(i ...interface{}) string {
	return fmt.Sprint(i...)
}

func (f *FixedLine) SplitSized(leftText, bodyText, rightText string) (string, string, string) {
	if f.haveSize {
		rightSize := f.RPanelWidth
		leftSize := f.LPanelWidth
		// step by step reduction of right and left if they need to much space
		if rightSize+leftSize > f.width {
			rightSize = 0
			rightText = ""
			if leftSize > f.width && !f.LPanelMustShown {
				leftSize = 0
				leftText = ""
			} else {
				leftText = stringPadLeft(leftText, leftSize)
			}
		} else {
			rightText = stringPadLeft(rightText, rightSize)
			leftText = stringPadLeft(leftText, leftSize)
		}
		maxSize := f.width - (rightSize + leftSize + f.BufferWith)
		bodyText = stringPadRight(bodyText, maxSize)
	}
	return leftText, bodyText, rightText
}

func stringPadLeft(str string, max int) string {
	if max < 1 {
		return ""
	}
	if len(str) > max {
		str = str[:max]
	}
	return fmt.Sprintf("%"+fmt.Sprint(max)+"s", str)
}

func stringPadRight(str string, max int) string {
	if max < 1 {
		return ""
	}
	if len(str) > max {
		str = str[:max]
	}
	return fmt.Sprintf("%-"+fmt.Sprint(max)+"s", str)
}
