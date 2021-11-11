package main

import (
	"bufio"
	"fmt"
	"gocryptotrader/common"
	"gocryptotrader/currency"
	"gocryptotrader/workerpool"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

const (
	baseDir = "/home/bijan/work/crypto/gateiodata"
	baseCmd = "dbseed candle file --exchange %s --base %s --quote %s --interval 60 --asset spot --filename %s"
)

func main() {
	log.SetFlags(log.Ltime)

	// For monitoring purpose.
	// waitC := make(chan bool)
	// go func() {
	// 	for {
	// 		log.Printf("[main] Total current goroutine: %d", runtime.NumGoroutine())
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	// Start Worker Pool.
	totalWorker := 1
	wp := workerpool.NewWorkerPool(totalWorker)
	wp.Run()

	type result struct {
		id    int
		value string
	}

	totalTask := 0
	resultC := make(chan result, 10000)

	finished := finishedSymbols()
	symbolDirs, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	// fmt.Println("finished", finished)
	for _, symbolDir := range symbolDirs {
		if symbolDir.IsDir() {
			if !inSymbolList(symbolDir.Name()) {
				continue
			}
			// fmt.Println("checking", symbolDir.Name())
			files, _ := ioutil.ReadDir(path.Join(baseDir, symbolDir.Name()))
			for _, f := range files {
				// fmt.Println(f.Name())
				name := f.Name()
				if !strings.HasSuffix(name, ".csv") {
					continue
				}
				if isFinished(name, finished) {
					// fmt.Println("skipping", name)
					continue
				}

				fmt.Println("create task NOT IN FINISHED LIST", name)

				i += 1
				id := i + 1
				totalTask = i
				wp.AddTask(func() {
					log.Printf("[main] Starting task %s", name)
					task(name)
					resultC <- result{id, name}
					fmt.Println("finished", name)
				})
				fmt.Println("done adding task")
			}
		}
	}

	fmt.Println("waiting for", totalTask)
	for i := 0; i < totalTask; i++ {
		// <-resultC
		res := <-resultC
		log.Printf("[main] Task %s has been finished", res.value)
	}

	// <-waitC
}

func isFinished(filename string, finished []string) bool {
	for _, fin := range finished {
		if strings.EqualFold(fin, filename) {
			return true
		}
	}

	return false
}

func inSymbolList(dirname string) bool {
	for _, s := range common.Symbols() {
		if s == dirname {
			return true
		}
	}
	return false
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func task(fileName string) {
	fmt.Println("run task", fileName)
	dirName := strings.Split(fileName, "-")[0]
	c, err := currency.NewPairFromString(dirName)
	if err != nil {
		fmt.Println("cant find currency", err)
	}
	cmd := fmt.Sprintf(baseCmd, "gateio", c.Base.String(), c.Quote.String(), path.Join(baseDir, dirName, fileName))
	// fmt.Println(cmd)
	// out, err := exec.Command("bash", "-c", cmd).Output()
	command := exec.Command("bash", "-c", cmd)
	var waitStatus syscall.WaitStatus
	if err := command.Run(); err != nil {
		printError(err)
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
		fmt.Println("err running command", cmd, err)
		os.Exit(2)
	}
	// else {
	// 	// Command was successful
	// 	waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
	// 	printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	// }
	// fmt.Println("out", out)
	markFileFinished(fileName)
}

func markFileFinished(symbol string) {
	f, err := os.OpenFile("./finished.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	output := fmt.Sprintf("%s\n", symbol)
	if _, err := f.WriteString(output); err != nil {
		log.Println(err)
	}
}

func finishedSymbols() []string {
	file, err := os.Open("./finished.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pairs := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pairs = append(pairs, scanner.Text())
	}
	fmt.Println("loaded", len(pairs))
	return pairs
}
