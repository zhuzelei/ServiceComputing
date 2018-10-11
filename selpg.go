package main

//import packages
import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

//var something needed
var (
	startPage, endPage, lineNumber int
	f_Flag, readFile_Flag, lp_Flag bool
	desName, fileName              string
	lp_Control                     *exec.Cmd
	toPipe                         io.WriteCloser
	pipeErr                        error
)

func parseArgs() {
	//输入参数
	//--------
	//“-s Snum” start from page Snum
	flag.IntVar(&startPage, "s", -1, "Input Start Page ( >=  0)")
	//“-e Enum” end at page Enum
	flag.IntVar(&endPage, "e", -1, "Input End Page (Not less than Start Page")

	//-------
	//“-l num” or “-f” ?
	//“-l num” define how many lines one page contains / default is 72
	flag.IntVar(&lineNumber, "l", 72, "Input How many lines per page / default is 72")
	//“-f” define lineNumber By \f (ASCII value = 12)
	flag.BoolVar(&f_Flag, "f", false, "define How many lines per page By '\\f' (using '-f' wil ignore '-l num')")
	//“-d Des” Output to Destination
	flag.StringVar(&desName, "d", "", "Ensure output Destination")

	//------
	flag.Parse()
	// —d
	if desName != "" {
		lp_Flag = true
		//管道连接
		lp_Control = exec.Command("lp", "-d", desName)
		//Test (using 'cat')
		//lp_Control = exec.Command("cat")
		//指定管道输出至该程序的标准输出
		lp_Control.Stdout = os.Stdout
		lp_Control.Stderr = os.Stderr
		//向管道的输入为程序控制的输出
		toPipe, pipeErr = lp_Control.StdinPipe()
		lp_Control.Start()
	}

	//Error report
	if startPage < 0 {
		log.Fatalln(errors.New("Start Page Number should be greater than 0"))
	}
	if endPage < 0 {
		log.Fatalln(errors.New("End Page Number should be greater than 0"))
	}
	if startPage > endPage {
		log.Fatalln(errors.New("End Page Numebr should be greater than Start Page Number "))
	}

	if lineNumber != 72 && f_Flag == true {
		log.Printf("ps: Using '-f' will ignoring '-l num' %d\n", lineNumber)
	}

	//stream
	if len(flag.Args()) == 0 {
		readFile_Flag = false
	} else {
		if len(flag.Args()) > 1 {
			log.Println("Only accept one input stream !")
		}
		readFile_Flag = true
		fileName = os.ExpandEnv(flag.Args()[0])
		pwd, err := os.Getwd()
		if err != nil {
			fileName = pwd + fileName
		}
	}
}

//selpg work function
func work() {
	//reader
	var reader *bufio.Reader
	if readFile_Flag == true {
		inputFile, inputErr := os.Open(fileName)
		if inputErr != nil {
			log.Fatal("Cannot open this file ! \n")
		}
		defer inputFile.Close()
		reader = bufio.NewReader(inputFile)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}
	//processOutput()
	if f_Flag {
		pagectr := 1
		for {
			pChar, _, err := reader.ReadRune()
			if err == io.EOF {
				if lp_Flag {
					toPipe.Close()
					lp_Control.Wait()
				}
				break
			} else if err != nil {
				panic(err)
			}
			if pChar == '\f' {
				pagectr++
			}
			if pagectr >= startPage && pagectr <= endPage {
				if lp_Flag {
					toPipe.Write([]byte(string(pChar)))
				} else {
					fmt.Print(string(pChar))
				}
			}
		}
	} else {
		pagectr := 1
		linectr := 0
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				if lp_Flag {
					toPipe.Close()
					lp_Control.Wait()
				}
				break
			} else if err != nil {
				panic(err)
			}
			linectr++
			if linectr > lineNumber {
				pagectr++
				linectr = 1
			}
			if pagectr >= startPage && pagectr <= endPage {
				if lp_Flag {
					toPipe.Write([]byte(line))
				} else {
					fmt.Print(line)
				}
			}
		}
	}
}

func main() {
	parseArgs()
	work()
}
