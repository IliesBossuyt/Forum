package main

import (
	server "Forum/server"
)

func main() {
	var forum server.Forum
	server.Router(&forum)
}