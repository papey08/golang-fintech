package main

import (
	"homework8/internal/adapters/adrepo"
	"homework8/internal/adapters/user_repo"
	"homework8/internal/app"
	"homework8/internal/ports/httpgin"
	"os"
)

func main() {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New(), user_repo.New()), os.Stdout)
	if err := server.Listen(); err != nil {
		panic(err)
	}
}
