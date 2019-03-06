package dictscanner

import (
	"errors"
	"fmt"
	"strings"
	"text/scanner"
)

type Translation struct {
	Finnish  []string
	English  []string
	Comments []string
}

func ParseLine(line string) (Translation, error) {

	var s scanner.Scanner

	s.Init(strings.NewReader(line))
	s.Whitespace ^= 1 << '<' // don't skip tabs and new lines

	var t Translation
	var x int
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '\n':
			break
		case '<':
			break
		case '|':
			for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
				str := s.TokenText()
				t.English = append(t.English, str)
				//fmt.Printf("%v\n", t.English)
			}
		default:
			x++
			str := s.TokenText()
			if str[0] == '<' {
				break
			}
			if str[len(str)-1] == '<' {
				break
			}
			if str[0] >= 'A' && str[0] <= 'Ö' && len(str) > 1 {
				t.Finnish = append(t.Finnish, str)
				//fmt.Printf("%v\n", t.Finnish)
			}
		}
	}

	fmt.Printf("%v\n", t)
	fmt.Printf("fin %s\n", t.Finnish)
	fmt.Printf("eng %s\n", t.English)
	//os.Exit(1)
	return t, nil
}

func ParseLineForFinnishPart(line string) (*Translation, error) {
	parts := strings.Split(line, " ")

	var t Translation
	for _, word := range parts {
		if word[0] == '|' {
			break
		}
		if word[0] == '<' {
			continue
		}
		//var http_tag_name_end string
		if ((word[0] >= 'A' && word[0] <= 'Ö') || word[0] == '<') && len(word) > 1 {
			var next_is_comment bool = false
			for {
				cut_end, is_start_tag, is_end_tag, content, tag := ParseHttpTags(word)
				if len(tag) == 0 {
					if len(content) == 0 && len(word) > 0 {
						t.Finnish = append(t.Finnish, word)
					}
					break
				}
				if is_start_tag == true {
					word = word[cut_end+1:]
					next_is_comment = true
					continue
				}
				if (len(tag) > 0 && is_end_tag == true) || next_is_comment == true {
					if len(content) > 0 && next_is_comment == false {
						t.Finnish = append(t.Finnish, content)
						if len(tag) > 0 {
							t.Comments = append(t.Comments, tag)
						}
					} else {
						t.Comments = append(t.Comments, tag)
						if len(content) > 0 {
							t.Comments = append(t.Comments, ","+content)
						}
					}
				}
				if is_end_tag == true {
					//t.Finnish = append(t.Finnish, word[0:cut_end-(len(tag)+2)])
					next_is_comment = false
					word = word[cut_end+1:]
					if len(word) == 0 {
						break
					}
				} else {
					word = word[cut_end+1:]
					next_is_comment = true
				}
			}
		}
	}

	if len(t.Finnish) == 0 {
		return nil, errors.New("no Finnish words found")
	}

	var found bool = false
	for _, word := range parts {
		//word := parts[x]
		//fmt.Println(parts[i])
		if len(word) > 0 && word[0] == '|' {
			found = true
			continue
		}
		if found == true && len(word) > 0 {
			t.English = append(t.English, word)
		}
	}
	return &t, nil
}

func ParseHttpTags(word string) (int, bool, bool, string, string) {
	var tag_start int = -1
	var tag_end int = 0
	var content_end int = 0
	for i, ch := range word {
		//fmt.Printf("%c ", ch)
		if ch == '<' {
			tag_start = i
		}
		if ch == '>' && tag_start != -1 {
			tag_end = i
			break
		}
	}

	if (tag_end > tag_start) && tag_start > 0 {
		content_end = tag_start
	}
	if tag_start == -1 {
		return content_end, false, false, word[0:content_end], ""
	}
	var is_end_tag bool = false
	if word[tag_start+1] == '/' || word[tag_end-1] == '/' {
		//fmt.Printf("tag %s\n", word[tag_start:tag_end])
		is_end_tag = true
		if word[tag_start+1] == '/' {
			return tag_end, false, is_end_tag, word[0:content_end], word[tag_start+2 : tag_end]
		} else {
			return tag_end, false, is_end_tag, word[0:content_end], word[tag_start+1 : tag_end-1]
		}
	} else {
		var is_start_tag bool = false
		if tag_start != -1 && tag_end > tag_start && content_end == 0 {
			is_start_tag = true
		}
		//fmt.Printf("tag %s\n", word[tag_start:tag_end])
		return tag_end, is_start_tag, is_end_tag, word[0:content_end], word[tag_start+1 : tag_end]
	}
}
