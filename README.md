First you need to create module by using command `go mod init path/to/mod` then if you want to create a package you need to create a folder, the `.go` files in that folder are all in the same package with the name of that folder

To then import you can write `import "mod_name/package_name" `

#### To do:
- Wrap interactions with backends inside a struct
- Implement a Pool struct in the lb to handle also different 
types of routing policies for backends