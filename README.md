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
         -threads=10
```
             
### Create Users with PasswordHook

```
./okgo -command=createUsersWithHook
```

### Spinning Multiple threads

```
./okgo -command=createUsersWithHook -threads=3
```

