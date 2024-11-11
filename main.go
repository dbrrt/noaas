package main

import "dbrrt/noaas/routing"

func main() {
	r := routing.SetupServer()
	r.Run(":8080")
}
