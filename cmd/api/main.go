package main

func main() {
	app := NewApp()
	defer app.Exit()
	app.Run()
}
