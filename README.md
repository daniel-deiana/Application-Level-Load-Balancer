First you need to create module by using command `go mod init path/to/mod` then if you want to create a package you need to create a folder, the `.go` files in that folder are all in the same package with the name of that folder

To then import you can write `import "mod_name/package_name" `

#### To do:
- Add locking for operations on proxies slice
- Study better go packaging
- study better docker files


#### Reminder
Every time you execute `RUN` inside a docker file a new state and a new image is created so you can inspect that state by starting a container using the id provided by the docker build command, with that you can start a shell and check stuff 


