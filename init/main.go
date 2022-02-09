package main

import (
	i "github.com/rgraterol/where-is-my-movie/init/initializers"
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