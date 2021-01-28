package main

import (
	"strconv"
	"fmt"
	"time"
)

func createTestUsers(execConfig *ExecutorConfig, threads *int, okClient *OktaClient) {

	
	fmt.Println("Start time:", time.Now())

	// var failedRec int
	// var ratelimitRec int

		

		n := 0
		totalRec, _  := strconv.Atoi(execConfig.UserStatus)

	
		fmt.Println("login,email,firstName,lastName")
		
		
		for n < (totalRec){

			// ss := strconv.Itoa(n)

			rec := fmt.Sprintf("testUser%d@test.com,testUser%d@test.com,Test,User%d",n ,n, n)

			fmt.Println(rec)
		
		n++

	}
	fmt.Println("End time:", time.Now())
	fmt.Println("total records created::", totalRec)

}