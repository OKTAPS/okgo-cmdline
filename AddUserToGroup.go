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

func addUsersToGroupAll(execConfig *ExecutorConfig, threads *int, okClient *OktaClient) {

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
				addEachUserToGroup(record[0], execConfig.GroupId, okClient)
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

func addEachUserToGroup(User string, GroupId string, okClient *OktaClient) (jobStatus bool) {

	user, resp, err := okClient.Client.User.GetUser(*okClient.Ctx, strings.TrimSpace(User))

	if err != nil {

		if resp.StatusCode == 404 {
			//fmt.Println(err)
			fmt.Println("Invalid User::", strings.TrimSpace(User))
			return false
		}
	}

	//Add user to group

	if user.Status == "DEPROVISIONED" {
		fmt.Println("Deprovisioned User::", strings.TrimSpace(User))
		return false
	} else {

		resp, err = okClient.Client.Group.AddUserToGroup(*okClient.Ctx, GroupId, user.Id)

		if resp.StatusCode == 429 {
			fmt.Println("Ratelimit Error:", resp.StatusCode)
			return false
		}

		if err != nil {

			fmt.Println("Unknown Error::", (*user.Profile)["login"])
			return false
		}

		fmt.Println("Added User::", strings.TrimSpace(User))

	}
	// fmt.Println(user.Id)
	return true
}
