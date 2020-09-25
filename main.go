package main

import (
	"fmt"
)

func main() {
	ghc := GHClient{"araji", "testrepo"}
	slacker := Slacker{"araji", "password"}
	fmt.Printf("getting push requests\n")
	ghc.GetPushRequests()
	slacker.SendMessage(&SlackMessage{Channel: "mychannel", Text: "ALERT"})

}
