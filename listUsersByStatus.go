package main

import (
	"fmt"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func listUsersByStatusAll(execConfig *ExecutorConfig, threads *int, okClient *OktaClient) {

	fmt.Println("Start time:", time.Now())

	listUsersByStatus(okClient, execConfig)

	fmt.Println("End time:", time.Now())

	fmt.Println("Successfully Retrieved Records::")

}

func listUsersByStatus(okClient *OktaClient, execConfig *ExecutorConfig) {

	// filter := query.NewQueryParams(query.WithLimit(200))

	filter := query.NewQueryParams(query.WithFilter("status eq \""+execConfig.UserStatus+"\""), query.WithLimit(200))

	users, resp, err := okClient.Client.User.ListUsers(*okClient.Ctx, filter)

	if err != nil {
		if resp.StatusCode == 404 {
			fmt.Println("API Call failed")
		}

	}

	for _, user := range users {
		fmt.Println((*user.Profile)["login"], ",", user.Status)
	}

	if resp.HasNextPage() {
		NextPage(okClient, resp)
	}

}

func NextPage(okClient *OktaClient, resp *okta.Response) {

	var nextUserSet []*okta.User
	resp, _ = resp.Next(*okClient.Ctx, &nextUserSet)

	for _, user := range nextUserSet {
		fmt.Println((*user.Profile)["login"], ",", user.Status)
	}

	if resp.HasNextPage() {
		NextPage(okClient, resp)
	}

}
