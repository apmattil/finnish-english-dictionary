package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"finnish-english-dictionary"
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
		//var translations []Translation
		t, err := dictscanner.ParseLineForFinnishPart(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if len(t.English[0]) > 0 && len(t.Finnish[0]) > 0 {
			for i, word := range t.Finnish {
				if i > 0 {
					fw.WriteString(",")
				}
				fw.WriteString(word)
			}
			fw.WriteString("\t")
			for i, word := range t.English {
				if i > 0 {
					fw.WriteString(" ")
				}
				fw.WriteString(word)
			}
			fw.WriteString("\n")
		}
	}
	return nil
}

/* foobar*/
func HandleFile(f os.File, fw os.File) error {
	rd := bufio.NewReader(&f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return errors.New("read error")
		}
		fmt.Printf("line:%s\n", line)
		trans, perr := dictscanner.ParseLine(line)
		if perr != nil {
			return err
		}
		fmt.Printf("trans:%v\n", trans)

		for i, str := range trans.Finnish {
			if i > 0 {
				fw.WriteString(",")
			}
			fw.WriteString(str)
		}
		fw.WriteString("\t")
		for _, str := range trans.English {
			fw.WriteString(str)
			fw.WriteString(" ")
		}
		fw.WriteString("\r\n")
	}
	return nil
}

func hello() {
	fmt.Println("joo")
}
