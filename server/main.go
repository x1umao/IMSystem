package main

func main() {
	Server := NewServer("localhost",8080)
	Server.Start()
}
