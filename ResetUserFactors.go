package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func resetUserFactorsAll(execConfig *ExecutorConfig, threads *int, okClient *OktaClient) {

	fmt.Println("Start time:", time.Now())

	var successRec int
	// var failedRec int
	// var ratelimitRec int

	for _, f := range *execConfig.FilesToProcess {

		file, err := os.Open(f)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		i := 0

		fmt.Println("Processing Thread count:", *threads)
		fmt.Println("User", ",", "Status")
		ch := make(chan []string, *threads)

		var wg sync.WaitGroup

		reader := csv.NewReader(file)

		for {
			i++
			record, err := reader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}
			wg.Add(1)

			go func(record []string, m int) {
				defer wg.Done()
				ch <- record
				resetUserFactor(record[0], okClient, execConfig)
				time.Sleep(25 * time.Millisecond)
				<-(ch)
			}(record, i)

			successRec++
		}

		wg.Wait()
		close(ch)

		fmt.Println("End time:", time.Now())

		fmt.Println("Successfully Processed Records::", successRec)

	}

}

func resetUserFactor(User string, okClient *OktaClient, execConfig *ExecutorConfig) (jobStatus bool) {

	// fmt.Println("inside sub method")
	// fmt.Println(execConfig.GroupId)
	resp, err := okClient.Client.User.ResetFactors(*okClient.Ctx, strings.TrimSpace(User))

	if err != nil {

		if resp.StatusCode == 404 {
			if resp.StatusCode == 404 {
				if execConfig.UserStatus == "ALL" {
					fmt.Println(strings.TrimSpace(User), ",Invalid User")
				}
			}

			return false

		}

	}

	// fmt.Println(user.Id)

	return true

}
