package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	dictscanner "finnish-english-dictionary"
)

//  ./mobigen.exe utf8 finnish-english-dict.opf

func main() {
	f_opf, err := os.Create("fin-eng.opf")
	if err != nil {
		fmt.Printf("can not open %s\n", err.Error())
		panic(err)
	}
	printOpfHeader(f_opf)
	PrintOpfTailer(f_opf)

	// Open the file and scan it.
	f, err1 := os.Open("data2.adj")
	if err1 != nil {
		fmt.Printf("can not open\n")
		panic(err1)
	}

	fw, werr := os.Create("parsed.txt")
	if werr != nil {
		fmt.Printf("can not open %s\n", werr.Error())
		panic(werr)
	}

	fr, err := os.Open("out.txt")
	if err != nil {
		fmt.Printf("can not open\n")
		panic(err)
	}

	defer func() {
		f.Close()
		fw.Close()
		fr.Close()
		f_opf.Close()
	}()

	err = ScanFile(f, fw, fr)
	if err != nil {
		fmt.Printf("handle failed %s\n", err.Error())
		panic(err)
	}
}

func ScanFile(f *os.File, fw *os.File, fr *os.File) error {
	scanner := bufio.NewScanner(f)

	var translations []dictscanner.Translation
	for scanner.Scan() {
		line := scanner.Text()
		t, err := dictscanner.ParseLineWords(line, fr)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		translations = append(translations, *t)
	}

	var lines []string
	for i, _ := range translations {
		translations[i].TransformToLines(&lines)
	}
	sortByFinnishAndLen(lines)
	/*
		for _, line := range lines {
			fmt.Printf("'%s'\n", line)
		}
	*/

	handleDublicates(&lines)

	for _, line := range lines {
		if len(line) > 0 {
			if line[0] != '-' {
				fw.WriteString(line)
				fw.WriteString("\r\n")
			}
		}
	}
	for _, line := range lines {
		if len(line) > 0 {
			if line[0] == '-' {
				fw.WriteString(line)
				fw.WriteString("\r\n")
			} else {
				break
			}
		}
	}

	f_html, werr2 := os.Create("out0.html")
	if werr2 != nil {
		fmt.Printf("can not open\n")
		panic(werr2)
	}
	writeHtmlPageHead(f_html)
	for _, line := range lines {
		if len(line) > 0 {
			if line[0] != '-' {
				writeHtmlTag(f_html, line)
			}
		}
	}
	writeHtmlTail(f_html)
	f_html.Close()
	return nil
}

func sortByFinnishAndLen(lines []string) {
	sort.SliceStable(lines, func(i, j int) bool {
		parts1 := strings.Split(lines[i], "\t")
		parts2 := strings.Split(lines[j], "\t")
		//mi, mj := parts1[0], parts2[0]
		switch {
		case parts1[0] != parts2[0]:
			return parts1[0] < parts2[0]
		default:
			return parts1[0] < parts2[0]
		}
	})
}

func handleDublicates(lines *[]string) {
	i := 0
	for j := 1; j < len(*lines); j++ {
		if i == j {
			continue
		}

		parts1 := strings.Split((*lines)[i], "\t")
		parts2 := strings.Split((*lines)[j], "\t")

		if len(parts1[0]) > 0 && len(parts2[0]) > 0 {
			if parts1[0] == parts2[0] {
				if parts1[1] != parts2[1] {
					(*lines)[j] = (*lines)[j] + "; " + parts1[1]
					(*lines)[i] = ""
				} else {
					(*lines)[i] = ""
				}
			}
		}
		i++
	}
}

func writeHtmlTag(f *os.File, line string) {
	parts := strings.Split(line, "\t")
	var fin_re = regexp.MustCompile(`_`)
	fin_s := fin_re.ReplaceAllString(parts[0], ` `)
	var re = regexp.MustCompile(`;`)
	s := re.ReplaceAllString(parts[1], `</p>`+"\r\n"+`<p>`)
	re = regexp.MustCompile(`<p></p>`)
	s = re.ReplaceAllString(s, ``)
	f.WriteString(`<mbp:pagebreak/>`)
	f.WriteString(`<idx:entry name="word" scriptable="yes">` + "\r\n" +
		`<h2>` + "\r\n" +
		"\t" + `<idx:orth>` + fin_s + `</idx:orth><idx:key key="` + fin_s + `">` + "\r\n" +
		`</h2>` + "\r\n" +
		"<p>" + s + "</p>" + "\r\n" +
		`</idx:entry>`)
}

func writeHtmlPageHead(f *os.File) {
	f.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<html xmlns:idx="www.mobipocket.com" xmlns:mbp="www.mobipocket.com" xmlns:xlink="http://www.w3.org/1999/xlink">
  <body>
    <mbp:pagebreak/>
    <mbp:frameset>
      <mbp:slave-frame display="bottom" device="all" breadth="auto" leftmargin="0" rightmargin="0" bottommargin="0" topmargin="0">
        <div align="center" bgcolor="yellow"/>
        <a onclick="index_search()">Dictionary Search</a>
        </div>
      </mbp:slave-frame>` + "\r\n")
}

func writeHtmlTail(f *os.File) {
	f.WriteString(`</mbp:frameset>
              </body>
            </html>
            `)
}

func printOpfHeader(f *os.File) {
	f.WriteString(`<?xml version="1.0"?><!DOCTYPE package SYSTEM "oeb1.ent">

<!-- the command line instruction 'prcgen dictionary.opf' will produce the dictionary.prc file in the same folder-->
<!-- the command line instruction 'mobigen dictionary.opf' will produce the dictionary.mobi file in the same folder-->

<package unique-identifier="uid" xmlns:dc="Dublin Core">

<metadata>
	<dc-metadata>
		<dc:Identifier id="uid">fin-eng-dictionary</dc:Identifier>
		<!-- Title of the document -->
		<dc:Title><h2>Finnish to English dictionary</h2></dc:Title>
		<dc:Language>FI</dc:Language>
	</dc-metadata>
	<x-metadata>
	        <output encoding="utf-8" flatten-dynamic-dir="yes"/>
		<DictionaryInLanguage>FI</DictionaryInLanguage>
		<DictionaryOutLanguage>EN</DictionaryOutLanguage>
	</x-metadata>
</metadata>

<!-- list of all the files needed to produce the .prc file -->
<manifest>`)
}

func PrintOpfTailer(f *os.File) {
	f.WriteString(`<!-- list of all the files needed to produce the .prc file -->
<manifest>
  <item href="en-fi-cover.jpg" id="my-cover-image" media-type="image/jpeg"/>
 <item id="dictionary0" href="en-fi0.html" media-type="text/x-oeb1-document"/>
</manifest>


<!-- list of the html files in the correct order  -->
<spine>
	<itemref idref="dictionary0"/>
</spine>

<tours/>
<guide> <reference type="search" title="Dictionary Search" onclick= "index_search()"/> </guide>
</package>
`)
}

/*
lisa:
           <idx:orth>journal
           <idx:infl>
             <idx:iform value="journals"/>
           </idx:infl>
         </idx:orth>
orig:
           <idx:orth>
					 journal
					 </idx:orth>
					 <idx:key key="journal">

*/
