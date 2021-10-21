package engine

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/candle"
	"gocryptotrader/wpool"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

const (
	defaultDataImporterBaseDir     = "/home/bijan/work/crypto/kraken_data"
	defaultDataImporterBaseCmd     = "dbseed candle file --exchange %s --base %s --quote %s --interval 60 --asset spot --filename %s"
	finishedFilename               = "finished.log"
	defaultDataImporterWorkercount = 12
)

var (
	errDefault = errors.New("wrong argument type")
	execFn     = func(ctx context.Context, args interface{}) (interface{}, error) {
		argVal, ok := args.(int)
		if !ok {
			return nil, errDefault
		}

		return argVal * 2, nil
	}
)

type DataImporter struct {
	baseDir     string
	baseCmd     string
	workerCount int
}

func SetupDataImporter() *DataImporter {
	return &DataImporter{
		baseDir:     defaultDataImporterBaseDir,
		baseCmd:     defaultDataImporterBaseCmd,
		workerCount: defaultDataImporterWorkercount,
	}
}

func (d *DataImporter) Run(exchange string) {
	log.SetFlags(log.Ltime)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	fmt.Println("running", d.workerCount, "workers")
	wp := wpool.New(d.workerCount)
	go wp.GenerateFrom(d.krakenJob())
	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			_, err := strconv.ParseInt(string(r.Descriptor.ID), 10, 64)
			if err != nil {
				log.Fatalf("unexpected error: %v", err)
			}
			val := r.Value.(string)
			fmt.Println("finished", val)
			exCount, _ := candle.CountExchange("kraken", 60, "spot")
			fmt.Println("exchange now has", (exCount / (1000)), "k candles")
			// if val != int(i)*2 {
			// 	log.Fatalf("wrong value %v; expected %v", val, int(i)*2)
			// }
		case <-wp.Done:
			fmt.Println("calling wp.done")
			return
		default:
		}
	}
	fmt.Println("done")
}

func printCommand(cmd *exec.Cmd) {
	fmt.Println(strings.Join(cmd.Args, " "))
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

// func task(fileName string) {
func (d *DataImporter) task(ctx context.Context, args interface{}) (interface{}, error) {
	fileName := args.(string)
	pairName := strings.Split(fileName, "_")[0]
	c, err := currency.NewPairFromString(pairName)
	if err != nil {
		fmt.Println("cant find currency", err)
	}

	// os.Exit(123)

	// if !lastCandle.Timestamp.IsZero() {
	// 	// fmt.Println("already have data for", c)
	// 	return fileName, nil
	// }

	cmd := fmt.Sprintf(d.baseCmd, "kraken", c.Base.String(), c.Quote.String(), path.Join(d.baseDir, fileName))
	command := exec.Command("bash", "-c", cmd)
	printCommand(command)

	var waitStatus syscall.WaitStatus
	if _, err := command.Output(); err != nil {
		fmt.Println("error")
		printError(err)
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
		fmt.Println("err running command", cmd, err)
		os.Exit(2)
	} else {
		// fmt.Println("success")
		// Command was successful
		waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
		// printOutput(output)
		// check results
		lastCandle, err := candle.Last("kraken",
			c.Base.String(),
			c.Quote.String(),
			60,
			"spot")
		if err != nil {
			fmt.Println("error getting last candle", err)
		}
		if lastCandle.Timestamp.IsZero() {
			fmt.Println("did not update correctly")
			os.Exit(123)
		} else {
			fmt.Println("last candle for", c, "is", lastCandle.Timestamp)
		}
	}
	return fileName, nil
}

func (d *DataImporter) krakenJob() []wpool.Job {
	jobsCount := 0
	jobs := make([]wpool.Job, 0)

	files, err := ioutil.ReadDir(d.baseDir)
	if err != nil {
		log.Fatal(err)
	}

	for i, f := range files {

		if strings.HasSuffix(f.Name(), "_1.csv") {
			// if !strings.EqualFold(f.Name(), "SRMGBP_1.csv") {
			// 	continue
			// }

			filePath := path.Join(d.baseDir, f.Name())
			fileCount, _ := lineCounter(filePath)
			fileName := f.Name()
			pairName := strings.Split(fileName, "_")[0]
			c, err := currency.NewPairFromString(pairName)
			if err != nil {
				fmt.Println("cant find currency", err)
			}
			dbCount, err := candle.Count("kraken",
				c.Base.String(),
				c.Quote.String(),
				60,
				"spot")

			// fmt.Println("lines in file", fileCount, "db", dbCount)
			if fileCount == 0 {
				os.Exit(123)
			}
			if int(dbCount) >= fileCount {
				fmt.Printf(".")
				// fmt.Println("good", c, dbCount)
				continue
			} else {
				fmt.Println(c, "has only", dbCount, "bars out of", fileCount)
			}
			// time.Sleep(time.Second * 5)

			jobs = append(jobs, wpool.Job{
				Descriptor: wpool.JobDescriptor{
					ID:       wpool.JobID(fmt.Sprintf("%v", i)),
					JType:    "anyType",
					Metadata: nil,
				},
				ExecFn: d.task,
				Args:   f.Name(),
			})
			jobsCount += 1
		}
	}
	fmt.Println("returned", len(jobs), "jobs")

	return jobs
}

func lineCounter(filepath string) (int, error) {
	file, _ := os.Open(filepath)
	r := bufio.NewReader(file)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
