package main

import (
	"fmt"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig("", "")
	token := oauth1.NewToken("", "")
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Search Tweets
	tweets, _, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: "NYGovCuomo",
		Count:      100,
	})

	for _, tweet := range tweets {
		if strings.Contains(tweet.Text, "briefing") && strings.Contains(tweet.Text, "ET") {
			fmt.Println(tweet.Text)
			timeEndIndex := strings.Index(tweet.Text, "ET")
			fmt.Println(tweet.Text[timeEndIndex-8 : timeEndIndex+2])
			return
		}
	}
}
