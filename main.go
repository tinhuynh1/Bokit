package main

import (
	"booking-svc/pkg/utils"
	"fmt"
	"time"
)

func main() {
	accessToken, _, _ := utils.GenerateTokens("user", []byte("Ym9va2luZy10aWNrZXQ="), time.Hour*24, time.Hour*24*30)
	fmt.Println(accessToken)
}
