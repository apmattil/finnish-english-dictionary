package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	dictscanner "finnish-english-dictionary"
)

func main() {
	// Open the file and scan it.
	f, err := os.Open("data.adj")
	if err != nil {
		fmt.Printf("can not open\n")
		panic(err)
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

/*
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		finnWords := strings.Split(parts[1], ",")
		// loop trough finnish words
		for i, word := range finnWords {
			// loop trough translations
			for j, t := range translations {
				// loop trough finnish words of translation
				for k, w := range t.Finnish {
					if w == word {
						alreadyFound := false
						var w2 string
						z := 0
						for z, w2 = range t.EnglishWordTranslations {
							if w2 == word {
								alreadyFound = true
							}
						}
						if alreadyFound == false {
							fmt.Printf("found %d:%d:%d:%d  %s : %s\n", i, j, k, z, parts[0], word)
							t.EnglishWordTranslations = append(t.EnglishWordTranslations, parts[0])
						}
						break
					}
					break
				}
			}
		}
	}
*/

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
