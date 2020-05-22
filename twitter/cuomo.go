package twitter

import (
	"fmt"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// FindCuomoBriefingTime returns the time of New York Governor Cuomo's next briefing by searching twitter
func FindCuomoBriefingTime() string {
	config := oauth1.NewConfig(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_KEY_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	twitterClient := twitter.NewClient(httpClient)

	// Search Tweets
	tweets, _, _ := twitterClient.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: "NYGovCuomo",
		Count:      100,
	})
	fmt.Println("Tweets retrieved")

	for _, tweet := range tweets {
		if strings.Contains(tweet.Text, "briefing") && strings.Contains(tweet.Text, "ET") {
			timeEndIndex := strings.Index(tweet.Text, "ET")
			briefingTime := tweet.Text[timeEndIndex-8 : timeEndIndex+2]
			fmt.Println(tweet.Text)
			fmt.Println(briefingTime)

			return briefingTime
		}
	}

	return ""
}
