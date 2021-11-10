package analyze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocryptotrader/common/file"
	"gocryptotrader/log"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (w *PortfolioWeights) Save(fpath string) error {
	if fpath == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Could not get working directory. Error: %v.\n", err)
			os.Exit(1)
		}
		fpath = fmt.Sprintf(
			"portfolio_analysis_%v.json",
			time.Now().Format("2006-01-02-15-04-05"))
		fpath = filepath.Join(wd, "results/pf", fpath)
		fmt.Println("saving", fpath)
	}
	writer, err := file.Writer(fpath)
	defer func() {
		if writer != nil {
			err = writer.Close()
			if err != nil {
				log.Error(log.Global, err)
			}
		}
	}()
	payload, err := json.MarshalIndent(w.Strategies, "", " ")
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(payload))
	return err
}
