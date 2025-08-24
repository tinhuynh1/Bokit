package main

import (
	"booking-svc/pkg/utils"
	"fmt"
	"time"
)

func main() {
	accessToken, _, _ := utils.GenerateTokens("admin", []byte("Ym9va2luZy10aWNrZXQ="), time.Hour*24*15, time.Hour*24*30)
	fmt.Println(accessToken)
}
