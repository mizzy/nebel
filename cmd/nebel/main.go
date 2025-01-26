package main

import (
	"github.com/alecthomas/kong"
	"github.com/mizzy/nebel"
)

var CLI struct {
	New struct {
		Title string `arg:"" name:"title" help:"Title of the new post." type:"title"`
	} `cmd:"" help:"Create a new post."`
	Generate struct {
	} `cmd:"" help:"Generate files."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "new <title>":
		err := nebel.CreateNewPost(CLI.New.Title)
		if err != nil {
			panic(err)
		}
	case "generate":
		err := nebel.Generate()
		if err != nil {
			panic(err)
		}
	default:
		panic(ctx.Command())
	}
}
