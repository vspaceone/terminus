# terminus
Mitgliedbezogene Dienste für den vspace.one

## API
* `/` Shows API information
* `/get/user` Get user information
  * `uid` Pass the uid for selecting the user
  * `token` Authenticate with a token
* `/new/user` Create a new user
  * `uid` Pass the uid to be associated with the new user
  * `username` Pass the wished username
  * `fullname` Pass the users fullname
  * `password` Pass the wished password
* `/auth` Authenticate as a specific user and get a token for authenticating until the end of the session
  * `uid` Pass the uid for selecting the user
  * `password` Pass the users password to authenticate 