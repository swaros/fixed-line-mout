package fixedlinemout

import (
	"errors"
	"regexp"
	"strings"
)

type Markup struct {
	Name      string
	Value     interface{}
	Reference string
}

type MarkupEntry struct {
	Text       string
	Properties []Markup
	Parsed     string
}

type MarkupParser struct {
	HandleErrors bool
	Entries      []MarkupEntry
}

type MarkupRunner struct {
	Exec func(mk Markup, current string) string
}

type Runner struct {
	HandleErrors bool
	Runners      map[string]MarkupRunner
}

func NewRunner(handleErrors bool) Runner {
	return Runner{Runners: make(map[string]MarkupRunner)}
}

func (mp *MarkupParser) ParseAll(runner Runner) (string, error) {
	current := ""
	mp.HandleErrors = runner.HandleErrors
	for _, me := range mp.Entries {
		if cur, err := me.ExecParse(runner, current, mp.HandleErrors); err != nil {
			if mp.HandleErrors {
				return cur, err
			}
		} else {
			current = cur
		}

	}
	return current, nil
}

func (mark *MarkupEntry) ExecParse(runner Runner, current string, stopOnErrors bool) (string, error) {
	for _, mup := range mark.Properties {

		if runner, exists := runner.Runners[mup.Name]; exists {
			current = runner.Exec(mup, current)
		} else {
			if stopOnErrors {
				return "", errors.New("undefined propertie: " + mup.Name)
			}
		}
	}
	return current, nil
}

func getMarkups(orig string, ref string) []Markup {
	var marks []Markup
	orig = strings.TrimLeft(orig, "<")
	orig = strings.TrimRight(orig, ">")
	parts := strings.Split(orig, ";")
	for _, pstr := range parts {
		pParts := strings.Split(pstr, "=")
		var m Markup
		m.Name = pParts[0]
		m.Reference = ref
		if len(pParts) > 1 {
			m.Value = pParts[1]
		}
		marks = append(marks, m)
	}
	return marks
}

// ParseMarkup is parsing the origin string
// and returns MarkupParser
func ParseMarkup(orig string) MarkupParser {
	var parsed MarkupParser
	if markups, found := splitByMarks(orig); found {
		allMk := len(markups)
		for mkIndex, mk := range markups {
			if indxStart := strings.Index(orig, mk); indxStart >= 0 {
				offset := len(mk)
				text := ""
				if mkIndex+1 < allMk {
					if end := strings.Index(orig[indxStart+offset:], "<"); end >= 0 {
						text = orig[indxStart+offset : indxStart+end+offset]
					}
				} else {
					text = orig[indxStart+offset:]
				}
				var entry MarkupEntry = MarkupEntry{
					Text:       text,
					Properties: getMarkups(mk, text),
				}
				parsed.Entries = append(parsed.Entries, entry)
			}
		}
	}

	return parsed
}

// splitByMarks extract all markups parts
func splitByMarks(orig string) ([]string, bool) {
	var result []string
	found := false
	re := regexp.MustCompile(`<[^>]+>`)
	newStrs := re.FindAllString(orig, -1)
	for _, s := range newStrs {
		found = true
		result = append(result, s)

	}
	return result, found
}
