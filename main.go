package main

func main() {

	c := InitCrawler()
	// listen and serve on 0.0.0.0:8080
	c.Router.Run()
}
