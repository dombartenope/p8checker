package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

type Request struct {
	Path   string
	Key    string
	Team   string
	Bundle string
}

func inputParser() Request {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Is there a p8 key located adjacent to this file? y/n")
	scanner.Scan()
	confirmation := scanner.Text()
	if confirmation != "y" && confirmation != "Y" {
		log.Fatalf("Please add the p8 file in the folder that this program is running and try again")
	}
	setPath := readLocal()

	fmt.Println("What is the Key ID associated with this p8 Key?")
	scanner.Scan()
	setKey := scanner.Text()

	fmt.Println("What is the Team ID associated with this p8 Key?")
	scanner.Scan()
	setTeam := scanner.Text()

	fmt.Println("What is the Bundle ID associated with this p8 Key?")
	scanner.Scan()
	setBundle := scanner.Text()

	req := Request{
		Path:   setPath,
		Key:    setKey,
		Team:   setTeam,
		Bundle: setBundle,
	}

	return req
}

func main() {

	req := inputParser()
	res := prod_req(req.Path, req.Key, req.Team, req.Bundle)

	switch res.StatusCode {
	case 400:
		if res.Reason == "DeviceTokenNotForTopic" {
			fmt.Printf("Everything matches but device token is not valid : \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		} else if res.Reason == "TopicDisallowed" {
			fmt.Printf("Check the Bundle ID, there might be a mismatch : \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		} else {
			fmt.Printf("Bad Request Response (This file is not invalid), please try uploading to the dashboard again: \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
			res = dev_req(req.Path, req.Key, req.Team, req.Bundle)
			fmt.Printf("Using Sandbox\n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		}
	case 403:
		fmt.Printf("There may be an incorrect Key/Team ID or p8 Token in use : \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	case 200:
		fmt.Printf("Successful Response : \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	default:
		fmt.Printf("Unknown Error, see APNs response: \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}

}

func prod_req(path, key, team, bundle string) *apns2.Response {
	authKey, err := token.AuthKeyFromFile(path)
	if err != nil {
		log.Fatal("token error:", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
		KeyID: key,
		// TeamID from developer account (View Account -> Membership)
		TeamID: team,
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = "b52c07acc8a575cbee2da8De21f3a0457292a07824040c7d8b3464bb5225059ab"
	notification.Topic = bundle
	// notification.Payload = []byte(`{"aps":{"alert":"test"}}`) // See Payload section below

	// If you want to test push notifications for builds running directly from XCode (Development), use
	// client := apns2.NewClient(cert).Development()
	// For apps published to the app store or installed as an ad-hoc distribution use Production()

	client := apns2.NewTokenClient(token)
	res, err := client.Push(notification)
	if err != nil {
		log.Fatal("Error:", err)
	}

	return res

}

func dev_req(path, key, team, bundle string) *apns2.Response {
	/*==========TODO==========*/
	//Make this take any p8 found in the cwd
	authKey, err := token.AuthKeyFromFile(path)
	if err != nil {
		log.Fatal("token error:", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
		KeyID: key,
		// TeamID from developer account (View Account -> Membership)
		TeamID: team,
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = "b52c07acc8a575cbee2da8De21f3a0457292a07824040c7d8b3464bb5225059ab"
	notification.Topic = bundle
	// notification.Payload = []byte(`{"aps":{"alert":"test"}}`) // See Payload section below
	client := apns2.NewTokenClient(token).Development()
	res, err := client.Push(notification)
	if err != nil {
		log.Fatal("Error:", err)
	}

	return res

}

func readLocal() string {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	var path string

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".p8" {
			path = file.Name()
			fmt.Println("Found p8 file")
		}
	}
	return path
}
