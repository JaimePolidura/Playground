package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	myUrl string
	delay int = 5
	waitGroup sync.WaitGroup
)

type myHtppResponse struct {
	response * http.Response
	err error
}

func main()  {
	if len(os.Args) == 1 {
		fmt.Println("Need a URL and a delay!")
		return
	}
	myUrl = os.Args[1]
	if len(os.Args) == 3 {
		t, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		delay = t
	}

	fmt.Println("Delay:", delay)
	c := context.Background()
	c, cancel := context.WithTimeout(c, time.Duration(delay)*time.Second)
	defer cancel()
	fmt.Printf("Connecting to %s \n", myUrl)
	waitGroup.Add(1)
	go connect(c)
	waitGroup.Wait()
	fmt.Println("Exiting...")
}


func connect(c context.Context) error {
	defer waitGroup.Done()

	data := make(chan myHtppResponse, 1)
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	req, _ := http.NewRequest("GET", myUrl, nil)

	go func() {
		response, err := httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
			data <- myHtppResponse{nil, err}
			return
		} else {
			pack := myHtppResponse{response, err}
			data <- pack
		}
	}()

	select {
		case <-c.Done():
			httpTransport.CancelRequest(req)
			<-data
			fmt.Println("The request was cancelled!")

			return c.Err()
		case ok := <-data:
			err := ok.err
			resp := ok.response
			if err != nil {
				fmt.Println("Error select:", err)
				return err
			}

			defer resp.Body.Close()
			realHTTPData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error select:", err)
				return err
			}
			fmt.Printf("Server Response: %s\n", realHTTPData)
	}

	return nil
}
