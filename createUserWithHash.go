package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func createUserWithHashAll(execConfig *ExecutorConfig, threads *int, okClient *OktaClient) {

	fmt.Println("Start time:", time.Now())

	var successRec int
	// var failedRec int
	// var ratelimitRec int

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
				createUserWithHash(record, okClient, execConfig, &firstRow)
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

func createUserWithHash(User []string, okClient *OktaClient, execConfig *ExecutorConfig, firstRow *[]string) (jobStatus bool) {

	// user, resp, err := okClient.Client.User.GetUser(*okClient.Ctx, strings.TrimSpace(User))

	h := &okta.PasswordCredentialHash{}

	p := &okta.PasswordCredential{
		Hash: h,
	}

	uc := &okta.UserCredentials{
		Password: p,
	}

	profile := okta.UserProfile{}

	for i, s := range *firstRow {

		if strings.TrimSpace(s) == "algorithm" {
			h.Algorithm = strings.TrimSpace(User[i])
		} else if strings.TrimSpace(s) == "workFactor" {
			wf, err := strconv.ParseInt(strings.TrimSpace(User[i]), 10, 64)
			if err != nil {
				return false
			}
			h.WorkerFactor = wf
		} else if strings.TrimSpace(s) == "salt" {
			h.Salt = strings.TrimSpace(User[i])
		} else if strings.TrimSpace(s) == "value" {
			h.Value = strings.TrimSpace(User[i])
		} else if strings.TrimSpace(s) == "saltOrder" {
			h.SaltOrder = strings.TrimSpace(User[i])
		} else {
			profile[s] = strings.TrimSpace(User[i])
		}

	}

	if len(h.Algorithm) == 0 || len(h.Value) == 0 {
		return false
	}

	u := &okta.CreateUserRequest{
		Credentials: uc,
		Profile:     &profile,
	}

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

		return false

	}

	if execConfig.UserStatus == user.Status || execConfig.UserStatus == "ALL" {
		fmt.Println((*user.Profile)["login"], ",", user.Status)
	}

	return true

}
