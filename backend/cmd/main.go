package main

import (
	"fmt"
	"theyudhiztira/oengage-backend/internal/server"
)

func main() {
	asciiArt := `
   ____  _______   ___________   ____________
  / __ \/ ____/ | / / ____/   | / ____/ ____/
 / / / / __/ /  |/ / / __/ /| |/ / __/ __/
/ /_/ / /___/ /|  / /_/ / ___ / /_/ / /___
\____/_____/_/ |_/\____/_/  |_\____/_____/

`
	fmt.Println(asciiArt)
	server.InitServer()
}
