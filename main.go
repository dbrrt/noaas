package main

func main() {
	res, _ := readRemoteUriPayload("https://pastebin.com/raw/v1qLcE15", false)
	print(res)
}
