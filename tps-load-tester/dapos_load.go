package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"time"

	"github.com/dispatchlabs/commons/types"
	"github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type Meter struct {
	ResultCount int
	Start       time.Time
	End         time.Time
	Total       int
}

// Private key storage (genisis key) for alpha testing
type Config struct {
	PrivateKey string
	From       string
	To         string
	DelegateIP string
}

func loadConfig(file string) Config {

	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func buildGRPCConnectionPool(connections int) *grpcpool.Pool {

	p, err := grpcpool.New(func() (*grpc.ClientConn, error) {
		return &grpc.ClientConn{}, nil
	}, 1, connections, 0)
	if err != nil {
		fmt.Println("The pool returned an error: %s", err.Error())
	}

	return p
}

func createSampleTransaction(cfg *Config) *types.Transaction {

	tx := types.NewTransaction(cfg.PrivateKey, 1,
		cfg.From,
		cfg.To, 1, time.Now())

	return tx
}

func sendTransaction(tx *types.Transaction, cfg *Config, mtr *Meter, pool *grpcpool.Pool) http.Response {

	byt := []byte(tx.String())
	buffer := new(bytes.Buffer)
	buffer.Write(byt)

	url := "http://" + cfg.DelegateIP + ":1975/v1/transactions"
	resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		fmt.Println(err.Error())
	} else {
	}
	return *resp
}

func run(cfg *Config, mtr *Meter, pool *grpcpool.Pool) {

	ret := createSampleTransaction(cfg)

	//fmt.Println(ret)

	resp := sendTransaction(ret, cfg, mtr, pool)

	if resp.StatusCode == 200 {
		if mtr.ResultCount < mtr.Total {
			mtr.ResultCount++
		} else {
			if mtr.End == mtr.Start {
				mtr.End = time.Now()

				fmt.Println("Calculating")
				fmt.Println("Total TX %d, Time Diff %d", mtr.Total, time.Since(mtr.Start))
				fmt.Println("DONE")
			}
		}
	}
}

func runLoad(cfg *Config, tx int, mtr *Meter, pool *grpcpool.Pool) {

	for i := 0; i <= tx; i++ {
		go run(cfg, mtr, pool)
	}

}

func wait() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...\n")
			//cleanup(services, c)
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func main() {

	var mtr Meter
	mtr.ResultCount = 0
	mtr.Start = time.Now()
	mtr.End = mtr.Start
	mtr.Total = 1000

	pool := buildGRPCConnectionPool(10)

	fmt.Println("Strating load test")
	cfg := loadConfig("./key.json")

	runLoad(&cfg, mtr.Total, &mtr, Pool)

	wait()
}
