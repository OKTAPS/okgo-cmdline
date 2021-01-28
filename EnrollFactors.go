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

func EnrollFactorsAll(execConfig *ExecutorConfig, threads *int, okClient *OktaClient) {

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
		fmt.Println("User", ",", "FactorType", "FactorStatus")
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
				enrollFactor(record, okClient, execConfig, &firstRow)
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

func enrollFactor(User []string, okClient *OktaClient, execConfig *ExecutorConfig, firstRow *[]string) (jobStatus bool) {

	var stat bool
	if execConfig.UserStatus == "ACTIVE" {
		stat = true
	} else {
		stat = false
	}

	query := query.NewQueryParams(query.WithActivate(stat), query.WithUpdatePhone(true))



	user, resp, err := okClient.Client.User.GetUser(*okClient.Ctx, strings.TrimSpace(User[0]))

	if err != nil {

		if resp.StatusCode == 404 {
			if resp.StatusCode == 404 {
				if execConfig.UserStatus == "ALL" {
					fmt.Println(strings.TrimSpace(User[0]), ",Invalid User")
				}
			}

			fmt.Println(strings.TrimSpace(User[0]), ",Invalid User", "404")
			return false

		}

	}

	userId := user.Id

	for i, s := range *firstRow {


		if (s == "sms" && len(strings.TrimSpace(User[i]))  > 0 ) {

			factorProfile := okta.NewSmsUserFactorProfile()

			factorProfile.PhoneNumber = User[i]

			factor := okta.NewSmsUserFactor()
			factor.Profile = factorProfile

			_, resp, err := okClient.Client.UserFactor.EnrollFactor(*okClient.Ctx, userId , factor, query)

			if err != nil {
		
				fmt.Println(err)
				//implement Loggin framework
			
			}
	

			fmt.Println(strings.TrimSpace(User[0]), ",sms,", resp.Status)
		} else if (s == "voice" && len(strings.TrimSpace(User[i]) ) > 0 ) {
			factorProfile := okta.NewCallUserFactorProfile()
			factorProfile.PhoneNumber = User[i]

			factor := okta.NewCallUserFactor()
			factor.Profile = factorProfile

			_, resp, _ := okClient.Client.UserFactor.EnrollFactor(*okClient.Ctx, userId, factor, query)

			if err != nil {
		
			//implement Loggin framework
		
			}

			fmt.Println(strings.TrimSpace(User[0]), ",voice,", resp.Status)
		} else if (s == "email" && len(strings.TrimSpace(User[i])) > 0  ){
			factorProfile := okta.NewEmailUserFactorProfile()
			factorProfile.Email = User[i]

			factor := okta.NewEmailUserFactor()
			factor.Profile = factorProfile

			_, resp, err := okClient.Client.UserFactor.EnrollFactor(*okClient.Ctx, userId , factor, query)

			if err != nil {
		
				//implement Loggin framework
			
			}
	

			fmt.Println(strings.TrimSpace(User[0]), ",email,", resp.Status )
		} 

	}

	return true

}
