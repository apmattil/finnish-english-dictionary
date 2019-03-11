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

	defer func() {
		f.Close()
		fw.Close()
	}()

	err = ScanFile(f, fw)
	if err != nil {
		fmt.Printf("handle failed %s\n", err.Error())
		panic(err)
	}
}

func ScanFile(f *os.File, fw *os.File) error {
	scanner := bufio.NewScanner(f)

	var translations []dictscanner.Translation
	for scanner.Scan() {
		line := scanner.Text()
		t, err := dictscanner.ParseLineWords(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		translations = append(translations, *t)
	}
	var lines []string
	for i, _ := range translations {
		translations[i].TransfomToLines(&lines)
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
	j := 0
	store_index := -1
	for i := 1; i < len(*lines); i++ {
		if i == j {
			continue
		}

		parts1 := strings.Split((*lines)[j], "\t")
		parts2 := strings.Split((*lines)[i], "\t")

		if len(parts1[0]) > 0 && len(parts2[0]) > 0 {
			if parts1[0] == parts2[0] {
				if store_index == -1 {
					store_index = j
					(*lines)[store_index] = (*lines)[j]
				}
				(*lines)[store_index] = (*lines)[store_index] + "; " + parts2[1]
				(*lines)[i] = ""
				if store_index != j {
					(*lines)[j] = ""
				}
				i++
				j++
			} else {
				store_index = -1
			}
		}
		j++
	}
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
