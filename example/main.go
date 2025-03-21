package main

import (
	"fmt"

	"github.com/HeavyBro1337/akevitt"
)

func main() {
	game := akevitt.NewGame(
		akevitt.Telnet(":1999"),
		akevitt.SSH(":1998", "id_rsa"),
	).
		Handle(func(ctx *akevitt.Context) {
			fmt.Fprintln(ctx, "Hello, World!")
		}).
		Plugin(akevitt.NewShellPlugin("# ", func(ctx *akevitt.Context, e *akevitt.Engine, s string) {
			fmt.Println("User says: " + s)

			for _, v := range e.Sessions() {
				if v == ctx {
					continue
				}

				fmt.Fprintf(v, "\nUser says: %s\n", s)
			}
		})).
		Engine()

	if err := game.Run(); err != nil {
		panic(err)
	}
}
