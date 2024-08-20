package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

type Request struct {
	Path   *string
	Key    *string
	Team   *string
	Bundle *string
	Token  *string
}

func flagParser() Request {
	setPath := flag.String("path", "", "desc: The path to the file that should be checked")
	setKey := flag.String("key", "", "desc: The Key ID")
	setTeam := flag.String("team", "", "desc: The Team ID")
	setBundle := flag.String("bundle", "", "desc: The Bundle ID")
	setPushToken := flag.String("pushtoken", "", "desc: The device Push Token")
	flag.Parse()

	req := Request{
		Path:   setPath,
		Key:    setKey,
		Team:   setTeam,
		Bundle: setBundle,
		Token:  setPushToken,
	}

	return req
}

func main() {
	tokenInfo := flagParser()
	//Example flag input :
	// go run main.go -path="my_p8.p8" -key="some_key" -team="some_key" -bundle="com.dom.support" -pushtoken="b52c07acc8a575cbee2da8e21f3a0457292a07824040c7d8b3464bb5225059ab"                                                                                                        ─╯

	/*==========TODO==========*/
	//Currently need to input everything manually using flags
	//Possibly add input listener for user to add data with terminal
	//Possibly read from txt or env file and input from there
	/*==========TODO==========*/

	authKey, err := token.AuthKeyFromFile(*tokenInfo.Path)
	if err != nil {
		log.Fatal("token error:", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
		KeyID: *tokenInfo.Key,
		// TeamID from developer account (View Account -> Membership)
		TeamID: *tokenInfo.Team,
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = *tokenInfo.Token
	notification.Topic = *tokenInfo.Bundle
	// notification.Payload = []byte(`{"aps":{"alert":"test"}}`) // See Payload section below

	// If you want to test push notifications for builds running directly from XCode (Development), use
	// client := apns2.NewClient(cert).Development()
	// For apps published to the app store or installed as an ad-hoc distribution use Production()

	client := apns2.NewTokenClient(token)
	res, err := client.Push(notification)
	if err != nil {
		log.Fatal("Error:", err)
	}

	switch res.StatusCode {
	case 400:
		if res.Reason == "DeviceTokenNotForTopic" {
			fmt.Printf("Everything matches but device token is not valid : \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		} else if res.Reason == "TopicDisallowed" {
			fmt.Printf("Check the Bundle ID, there might be a mismatch : \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		} else {
			fmt.Printf("Bad Request Response, but may not have a configuration error check Bundle and Push Token: \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		}
	case 403:
		fmt.Printf("There may be an incorrect Key/Team ID or p8 Token in use : \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	case 200:
		fmt.Printf("Successful Response : \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	default:
		fmt.Printf("Unknown Error, see APNs response: \n%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}

}
