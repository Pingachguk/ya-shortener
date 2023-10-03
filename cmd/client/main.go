package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func main() {
	endpoint := "http://localhost:8080/"

	fmt.Println("Введите длинный URL")
	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	long = strings.TrimSuffix(long, "\n")

	client := resty.New()
	req := client.R()
	req.Method = http.MethodPost
	req.URL = endpoint
	req.SetBody(long)
	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Статус-код ", res.StatusCode())
	fmt.Println(string(res.Body()))
}
