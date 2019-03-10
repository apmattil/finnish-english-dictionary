package main

import (
	"bufio"
	"fmt"
	"os"

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

	for scanner.Scan() {
		line := scanner.Text()
		var translations []dictscanner.Translation
		t, err := dictscanner.ParseLineWords(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		translations = append(translations, *t)
		for _, x := range translations {
			x.WriteTranslation(fw)
		}
	}
	return nil
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
