# okgo-cmdline

# Sample Config File
```
{
    "org_name": "subDomain",
    "base_url": "oktapreview.com",
    "api_token": "",
    "ignoreFirstRow": true
}
```

### Directory Structure

```
$ ls -lrt
./config
./input
okgo
```

## Commands to Run

## Usage

```
./okgo 
Commands: 
         -command=getUserId 
         -command=resetFactors 
         -command=listUsers <<STATUS>> 
         -command=enrollFactors <<STATUS>> 
         -command=createUserWithHash <<TARGET_STATUS>> 
         -command=createUser <<TARGET_STATUS>> 
         -command=createUsersWithHook 
         -command=createTestUsers <<USER_COUNT>> 
         -command=deleteUser 
         -command=getUserStatus <<FILTER_STATUS>> 
         -command=changetUserStatus <<TARGET_LIFECYCLE_STATUS>> <<Additional Query Params>> 
         -command=getUserNames 
         -command=addUsersToGroup <<GROUP_ID> 

 Threads: 
         -threads=<<number_of_threads>> (default set to 1)
```

##### Threads flag is optional and is defaulted to 1.

### Sample csv format for CreateUser/CreateUsersWithHook

firstName,lastName,email,login,<<attribute_variable_name>>
john,doe,john.doe@oktaice.com,john.doe@oktaice.com,<<attribute_value>>
             
### Create Users with PasswordHook

```
./okgo -command=createUsersWithHook
```

### Spinning Multiple threads

```
./okgo -command=createUsersWithHook -threads=3
```

### Creating Users with Active/STAGED status with multiple threads

```
./okgo -command=createUser -threads=3 STAGED

./okgo -command=createUser -threads=3 ACTIVE
```

### Sample csv format for enrollFactors (enrollFactors only supports sms,voice or email factors)
```
login,email,sms,voice
john.doe@oktaice.com,john.doe@oktaice.com,1234567890,1234567890
```
### Enroll in Users Factors 
```
./okgo -command=enrollFactors -threads=5 ACTIVE
```

### Reset Users Factors

```
./okgo -command=resetFactors -threads=3
```




