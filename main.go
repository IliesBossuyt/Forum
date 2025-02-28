package main

import (
	server "Forum/server"
)

func main() {
	server.Router(&server.Forum{})
}