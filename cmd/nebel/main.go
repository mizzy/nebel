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
	Serve struct {
		Addr string `default:"localhost:4000" help:"Address to listen on."`
	} `cmd:"" help:"Serve public/ and regenerate on changes."`
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
	case "serve":
		err := nebel.Serve(CLI.Serve.Addr)
		if err != nil {
			panic(err)
		}
	default:
		panic(ctx.Command())
	}
}
