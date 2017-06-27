package main

import (
	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {
	p.Task("darumasan", nil, func(c *do.Context) {
		c.Start("main.go")
	}).Src("*.go", "**/*.go")
}

func main() {
	do.Godo(tasks)
}
