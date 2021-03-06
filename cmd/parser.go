package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	dictscanner "finnish-english-dictionary"
)

/*
0. generate out.txt containing finnish to english tab separated translation file (see https://github.com/apmattil/finnish-english-dictionary)
1. strip license header from data files and copy to this directory
2. make title-page.html_pages_writen
3. make cover jpg
4. edit the opf header at printOpfHeader() to include cover image and page.
5. run with datafiles e.g  ./parser.exe data.adj data.adv data.noun data.verb
6. copy generated out*.html, fi-eng.opf, cover-image file and title-page.html to own directory e.g rel
7. get mobgen.exe ( http://web.archive.org/web/20070306012409/http://www.mobipocket.com/soft/prcgen/mobigen.zip)
8.  cd rel ; ./mobigen.exe utf8 fin-eng.opf
*/
//  ./mobigen.exe utf8 finnish-english-dict.opf

func main() {
	f_opf, err := os.Create("fin-eng.opf")
	if err != nil {
		fmt.Printf("can not open %s\n", err.Error())
		panic(err)
	}

	fw, werr := os.Create("parsed.txt")
	if werr != nil {
		fmt.Printf("can not open %s\n", werr.Error())
		panic(werr)
	}

	//fr, err := os.Open("out.txt")
	b, err := ioutil.ReadFile("out.txt")
	if err != nil {
		fmt.Printf("can not open\n")
		panic(err)
	}
	f_finn_translations := bufio.NewScanner(bytes.NewReader(b))

	defer func() {
		//f.Close()
		fw.Close()
		//fr.Close()
		f_opf.Close()
	}()

	scanFiles := os.Args[1:]
	html_pages_writen := 0
	var translations []dictscanner.Translation
	// Open the file and scan it.
	for i, file_name := range scanFiles {
		f, err1 := os.Open(file_name)
		if err1 != nil {
			fmt.Printf("can not open\n")
			panic(err1)
		}
		fmt.Printf("%d: handling %s\r\n", i, file_name)

		err = ScanFile(f, &translations, f_finn_translations, &html_pages_writen)
		if err != nil {
			fmt.Printf("scan file failed %s\n", err.Error())
			panic(err)
		}
	}

	err = HandleTranslations(translations, fw, &html_pages_writen)
	if err != nil {
		fmt.Printf("handle translations failed %s\n", err.Error())
		panic(err)
	}
	printOpfHeader(f_opf)
	PrintOpfTailer(f_opf, html_pages_writen)

}

func ScanFile(f *os.File, translations *[]dictscanner.Translation, f_finn_translations *bufio.Scanner, html_pages_writen *int) error {
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		t, err := dictscanner.ParseLineWords(line, f_finn_translations)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		*translations = append(*translations, *t)
	}
	return nil
}

func HandleTranslations(translations []dictscanner.Translation, fw *os.File, html_pages_writen *int) error {
	// TODO: make better middle format than lines
	var lines []string
	for i, _ := range translations {
		translations[i].TransformToLines(&lines)
	}
	fmt.Println("start shorting")
	sortByFinnishAndLen(lines)

	fmt.Println("start handle dublicates")
	handleDublicates(&lines)

	fmt.Println("write parsed.txt file")
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
	fmt.Println("write html files")
	WriteHtmlFiles(lines, html_pages_writen)
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

func WriteHtmlFiles(lines []string, pages_writen *int) {
	WriteHtmlFile(lines, false, pages_writen)
	//pages_writen = pages_writen + x
	//*pages_writen++
	WriteHtmlFile(lines, true, pages_writen)
}

func WriteHtmlFile(lines []string, write_under_scores bool, startId *int) {
	pages_writen := startId
	i := 0
	var f_html *os.File = nil
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '-' && write_under_scores == false {
			continue
		}
		if line[0] == '-' {
			if write_under_scores {
				if f_html == nil {
					f_html = createHtmlFile(pages_writen)
				}
				writeHtmlTag(f_html, line)
				i++
			}
		} else {
			if write_under_scores == false {
				if f_html == nil {
					f_html = createHtmlFile(pages_writen)
				}
				writeHtmlTag(f_html, line)
				i++
			}
		}
		if i == 999 {
			writeHtmlTail(f_html)
			f_html.Close()
			*pages_writen++
			i = 0
			f_html = nil
		}
	}
	if (i%999) != 0 && f_html != nil {
		writeHtmlTail(f_html)
		f_html.Close()
		*pages_writen++
	}
}

func createHtmlFile(pages_writen *int) *os.File {
	var err error
	var f_html *os.File = nil
	f_html, err = os.Create("out" + strconv.Itoa(*pages_writen) + ".html")
	if err != nil {
		fmt.Printf("can not open\n")
		panic(err)
	}
	fmt.Printf("writing %s\n", "out"+strconv.Itoa(*pages_writen)+".html")
	writeHtmlPageHead(f_html)
	return f_html
}

func writeHtmlTag(f *os.File, line string) {
	parts := strings.Split(line, "\t")
	var fin_re = regexp.MustCompile(`_`)
	fin_s := fin_re.ReplaceAllString(parts[0], ` `)
	var re = regexp.MustCompile(`;`)
	s := re.ReplaceAllString(parts[1], "\t\t"+`</p>`+"\r\n"+`<p>`)
	re = regexp.MustCompile(`<p></p>`)
	s = re.ReplaceAllString(s, ``)
	f.WriteString(`<mbp:pagebreak/>` + "\r\n" + "\r\n")
	f.WriteString(`<idx:entry name="word" scriptable="yes">` + "\r\n" +
		"\t" + `<h2>` + "\r\n" +
		"\t\t" + `<idx:orth>` + fin_s + `</idx:orth>` + "\r\n" +
		"\t\t" + `<idx:key key="` + fin_s + `">` + "\r\n" +
		"\t\t" + `</h2>` + "\r\n" +
		"\t\t" + "<p>" + s + "</p>" + "\r\n" +
		`</idx:entry>` + "\r\n")
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

func PrintOpfTailer(f *os.File, num_of_pages int) {
	f.WriteString(`<!-- list of all the files needed to produce the .prc file -->
<manifest>
  <item href="english-finnish-cover.jpg" id="my-cover-image" media-type="image/jpeg"/>` + "\r\n")

	f.WriteString("\t" + `<item id="title-page" href="title-page.html" media-type="text/x-oeb1-document"/>` + "\r\n")
	for i := 0; i < num_of_pages; i++ {
		f.WriteString("\t" + `<item id="dictionary` + strconv.Itoa(i) + `" href="out` + strconv.Itoa(i) + `.html" media-type="text/x-oeb1-document"/>` + "\r\n")
	}
	f.WriteString(`</manifest>` + "\r\n")

	f.WriteString(`<spine>` + "\r\n")
	f.WriteString("\t" + `<itemref idref="title-page"/>` + "\r\n")
	for j := 0; j < num_of_pages; j++ {
		f.WriteString("\t" + `<itemref idref="dictionary` + strconv.Itoa(j) + `"/>` + "\r\n")
	}
	f.WriteString(`</spine>` + "\r\n")

	f.WriteString(`
<guide> <reference type="search" title="Dictionary Search" onclick= "index_search()"/> </guide>
</package>
`)
	fmt.Printf("num_of_pages %d", num_of_pages)
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
