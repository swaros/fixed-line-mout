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
	LBracket     byte
	RBracket     byte
	Level1Sep    byte
	Level2Sep    byte
}

func NewRunner(handleErrors bool, bracketsAndSepByChar string) (Runner, error) {

	if len(bracketsAndSepByChar) != 4 {
		return Runner{}, errors.New("no valid bracket and seperatot definition found. you need to set these up like <>:;")
	}

	return Runner{
		Runners:   make(map[string]MarkupRunner),
		LBracket:  bracketsAndSepByChar[0],
		RBracket:  bracketsAndSepByChar[1],
		Level1Sep: bracketsAndSepByChar[2],
		Level2Sep: bracketsAndSepByChar[3],
	}, nil
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

func (r *Runner) getMarkups(orig string, ref string) []Markup {
	var marks []Markup
	orig = strings.TrimLeft(orig, string(r.LBracket))
	orig = strings.TrimRight(orig, string(r.RBracket))
	parts := strings.Split(orig, string(r.Level1Sep))
	for _, pstr := range parts {
		pParts := strings.Split(pstr, string(r.Level2Sep))
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
func (r *Runner) ParseMarkup(orig string) MarkupParser {
	var parsed MarkupParser
	if markups, found := r.splitByMarks(orig); found {
		allMk := len(markups)
		for mkIndex, mk := range markups {
			if indxStart := strings.Index(orig, mk); indxStart >= 0 {
				offset := len(mk)
				text := ""
				if mkIndex+1 < allMk {
					if end := strings.Index(orig[indxStart+offset:], string(r.LBracket)); end >= 0 {
						text = orig[indxStart+offset : indxStart+end+offset]
					}
				} else {
					text = orig[indxStart+offset:]
				}
				var entry MarkupEntry = MarkupEntry{
					Text:       text,
					Properties: r.getMarkups(mk, text),
				}
				parsed.Entries = append(parsed.Entries, entry)
			}
		}
	}

	return parsed
}

// splitByMarks extract all markups parts
func (r *Runner) splitByMarks(orig string) ([]string, bool) {
	var result []string
	found := false
	cmpStr := string(r.LBracket) + "[^" + string(r.RBracket) + "]+" + string(r.RBracket)
	re := regexp.MustCompile(cmpStr)
	newStrs := re.FindAllString(orig, -1)
	for _, s := range newStrs {
		found = true
		result = append(result, s)

	}
	return result, found
}
