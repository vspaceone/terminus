# terminus-server
Server for member specific services in a hackerspace

## API
* `/` Shows API information
* `/get/user` Get user information
  * `username`|`userid`|`uid` Pass the username/userid/uid for selecting the user
  * `token` Authenticate with a token
* `/new/user` Create a new user
  * `uid` Pass the uid to be associated with the new user
  * `username` Pass the wished username
  * `fullname` Pass the users fullname
  * `password` Pass the wished password
* `/auth` Authenticate as a specific user and get a token for authenticating until the end of the session
  * `userid` Pass the userid for selecting the user
  * `password` Pass the users password to authenticate 

----
`uid` : Is the unique tag id of the users nfc tag

## Configuration
Configuration is done by a json file, which needs to be in the current working directory.
At the moment this configuration file only specifies the database connection.
Sample configuration file (config.temp.json)
```
{
    "dbHost": "hostname",
    "dbName": "name of database",
    "dbUser": "name of db user",
    "dbPassw": "password of db user"
}
```

## Development
This software is developed under linux so it's not guaranteed to work on any other OS

* Be sure to have [golang installed and configured](https://golang.org/doc/code.html)
* `go get github.com/vspaceone/terminus-server`
* `go install github.com/vspaceone/terminus-server`
* `cd $GOHOME/src/github.com/vspaceone/terminus-server`
* Create a config.json according to the template for your needs
* to start `../../../../bin/terminus-server` from `$GOHOME/src/github.com/vspaceone/terminus-server` (currently config.json needs to be in CWD)