package main

import (
	i "github.com/rgraterol/movies-reviewers-wrapper/init/initializers"
)

func main() {
	run()
}

func run() {
	i.ConfigInitializer()
	i.LoggerInitializer()
	i.DatabaseInitializer()
	i.ServerInitializer()
}