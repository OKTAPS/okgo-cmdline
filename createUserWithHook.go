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

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func createUserWithHookAll(execConfig *ExecutorConfig, threads *int, okClient *OktaClient) {

	fmt.Println("Start time:", time.Now())

	var successRec int

	var firstRow []string

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

		//Read the First line
		// if *execConfig.config.IgnoreFirstRow == true {
		if firstRow, err = reader.Read(); err != nil {
			panic(err)
		}
		// }
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
				createUser(record, okClient, execConfig, &firstRow)
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

func createUserWithHook(User []string, okClient *OktaClient, execConfig *ExecutorConfig, firstRow *[]string) (jobStatus bool) {

	profile := okta.UserProfile{}

	for i, s := range *firstRow {

		profile[s] = strings.TrimSpace(User[i])

	}

	u := &okta.CreateUserRequest{
		Profile: &profile,
	}

	u.Credentials.Password.Hook.Type = "default"

	var stat bool
	if execConfig.UserStatus == "ACTIVE" {
		stat = true
	} else {
		stat = false
	}

	query := query.NewQueryParams(query.WithActivate(stat))

	user, resp, err := okClient.Client.User.CreateUser(*okClient.Ctx, *u, query)

	if err != nil {

		if resp.StatusCode == 404 {
			if execConfig.UserStatus == "ALL" {
				fmt.Println("Invalid User")
			}
		}

		fmt.Println(profile["login"], ",", "Error Creating User")

		return false

	}
	fmt.Println((*user.Profile)["login"], ",", user.Status)

	return true

}
