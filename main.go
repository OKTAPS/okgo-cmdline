package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type Configuration struct {
	Orgname        string `json:"org_name"`
	Baseurl        string `json:"base_url"`
	Apitoken       string `json:"api_token"`
	IgnoreFirstRow bool   `json:"ignoreFirstRow"`
}

type ExecutorConfig struct {
	InputAllArgs   []string
	GroupId        string
	config         *Configuration
	FilesToProcess *[]string
	UserStatus     string
}

type OktaClient struct {
	Ctx    *context.Context
	Client *okta.Client
}

func main() {

	command := flag.String("command", "help", "a string")

	threads := flag.Int("threads", 1, "an int")

	flag.Parse()

	execConf := ExecutorConfig{}

	execConf.InputAllArgs = flag.Args()

	setConfig(&execConf)

	filesToProcess, err := FilterDirsGlob("input", "*.csv")

	if err != nil {
		fmt.Println("No Files to Process")
	}

	execConf.FilesToProcess = &filesToProcess

	okClient, err := getOktaClient(&execConf)

	if err != nil {
		log.Fatal(err)
	}

	if *command == "getUserId" {
		fmt.Println("executing command")

		ggetUserIdsAll(&execConf, threads, okClient)

	} else if *command == "resetFactors" {

		fmt.Println("executing command")

		resetUserFactorsAll(&execConf, threads, okClient)
	} else if *command == "getUserNames" {

		fmt.Println("executing command")

		getUserNamesAll(&execConf, threads, okClient)
	} else if *command == "changetUserStatus" {

		if len(flag.Args()) == 0 {
			log.Fatal("missing target Status & , Please pass any of following status ACTIVATE, REACTIVATE, DEACTIVATE, SUSPEND, UNSUSPEND, DELETE, UNLOCK, RESET_PASSWORD, EXPIRE_PASSWORD, RESET_FACTORS, CLEAR_USER_SESSIONS")
		}

		execConf.UserStatus = flag.Args()[0]
		////TO DO
		// if execConfig.UserStatus = (ACTIVATE || REACTIVATE || DEACTIVATE || DELETE || RESET_PASSWORD)  {

		// }

		fmt.Println("executing command")

		changeUserStatusAll(&execConf, threads, okClient)

	} else if *command == "getUserStatus" {

		if len(flag.Args()) == 0 {
			log.Fatal("missing Status, Please pass any of following status ACTIVE, PROVISIONED, DEPROVISIONED, SUSPENDED,PASSWORD_EXPIRED, STAGED , ALL")
		}

		execConf.UserStatus = flag.Args()[0]

		fmt.Println("executing command")

		getUserStatusAll(&execConf, threads, okClient)
	} else if *command == "addUsersToGroup" {

		if len(flag.Args()) == 0 {
			log.Fatal("missing group ID, Please pass groupId")
		}

		execConf.GroupId = flag.Args()[0]

		addUsersToGroupAll(&execConf, threads, okClient)

	} else if *command == "createUserWithHash" {

		if len(flag.Args()) == 0 {
			log.Fatal("missing target Status, Please pass any of following status ACTIVE, STAGED")
		}

		execConf.UserStatus = flag.Args()[0]

		fmt.Println("executing command")

		createUserWithHashAll(&execConf, threads, okClient)

	} else if *command == "createUser" {

		if len(flag.Args()) == 0 {
			log.Fatal("missing target Status, Please pass any of following status ACTIVE, STAGED")
		}

		execConf.UserStatus = flag.Args()[0]

		fmt.Println("executing command")

		createUserAll(&execConf, threads, okClient)

	} else if *command == "createUsersWithHook" {

		if len(flag.Args()) == 0 {
			log.Fatal("missing target Status, Please pass any of following status ACTIVE, STAGED")
		}

		execConf.UserStatus = flag.Args()[0]

		fmt.Println("executing command")

		createUserWithHookAll(&execConf, threads, okClient)

	} else if *command == "deleteUser" {

		fmt.Println("executing command")

		deleteUserAll(&execConf, threads, okClient)

	} else if *command == "createTestUsers" {

		if len(flag.Args()) == 0 {
			log.Fatal("missing user count, please specify number of users to create")
		}

		execConf.UserStatus = flag.Args()[0]

		fmt.Println("executing command")

		createTestUsers(&execConf, threads, okClient)

	} else if *command == "listUsers" {
		if len(flag.Args()) == 0 {
			log.Fatal("missing Status, Please pass any of following status STAGED, PROVISIONED, ACTIVE, RECOVERY, PASSWORD_EXPIRED, LOCKED_OUT, DEPROVISIONED, SUSPENDED")
		}

		execConf.UserStatus = flag.Args()[0]

		fmt.Println("executing command")

		listUsersByStatusAll(&execConf, threads, okClient)

	} else if *command == "enrollFactors" {
		if len(flag.Args()) == 0 {
			log.Fatal("missing Status,  ACTIVE, INACTIVE")
		}

		execConf.UserStatus = flag.Args()[0]

		fmt.Println("executing command")

		EnrollFactorsAll(&execConf, threads, okClient)

	} else if *command == "help" {

		fmt.Println("Commands: \n \t -command=getUserId \n \t -command=resetFactors \n \t -command=listUsers <<STATUS>> \n\t -command=enrollFactors <<STATUS>> \n\t -command=createUserWithHash <<TARGET_STATUS>> \n\t -command=createUser <<TARGET_STATUS>> \n\t -command=createUsersWithHook \n\t -command=createTestUsers <<USER_COUNT>> \n\t -command=deleteUser \n\t -command=getUserStatus <<FILTER_STATUS>> \n \t -command=changetUserStatus <<TARGET_LIFECYCLE_STATUS>> <<Additional Query Params>> \n \t -command=getUserNames \n \t -command=addUsersToGroup <<GROUP_ID> \n\n Threads: \n \t -threads=10")

	} else {
		log.Fatal("Invalid command -- please run help for the list of commands")
	}

}

func setConfig(execConfig *ExecutorConfig) {

	configFile, err := os.Open("config/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	scConfig := Configuration{}
	err = decoder.Decode(&scConfig)
	if err != nil {
		fmt.Println("error:", err)
	}

	execConfig.config = &scConfig

}

// func FilterDirs(dir, suffix string) ([]string, error) {
// 	files, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res := []string{}
// 	for _, f := range files {
// 		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
// 			res = append(res, filepath.Join(dir, f.Name()))
// 		}
// 	}
// 	return res, nil
// }

func FilterDirsGlob(dir, suffix string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, suffix))
}

func getOktaClient(execConfig *ExecutorConfig) (*OktaClient, error) {

	ctx, client, err := okta.NewClient(context.TODO(), okta.WithOrgUrl("https://"+execConfig.config.Orgname+"."+execConfig.config.Baseurl), okta.WithToken(execConfig.config.Apitoken), okta.WithRateLimitMaxRetries(0), okta.WithRequestTimeout(60))

	if err != nil {
		fmt.Println(err)
	}

	okClient := OktaClient{}

	okClient.Client = client
	okClient.Ctx = &ctx

	return &okClient, err

}
