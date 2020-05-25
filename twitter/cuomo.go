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
		briefingTime := FindBriefingTimeInString(tweet.Text)
		if briefingTime != "" {
			return briefingTime
		}
	}

	return ""
}

// FindBriefingTimeInString returns the time of a briefing from a string
func FindBriefingTimeInString(tweet string) string {
	if !strings.Contains(tweet, "briefing") {
		return ""
	}

	fmt.Println(tweet)

	etIndex := strings.Index(tweet, "ET")
	if etIndex != -1 {
		briefingTime := tweet[etIndex-8 : etIndex+2]
		fmt.Println(briefingTime)
		return briefingTime
	}

	pmIndex := strings.Index(tweet, "PM")
	if pmIndex != -1 {
		briefingTime := tweet[pmIndex-5 : pmIndex+2]
		fmt.Println(briefingTime)
		return briefingTime
	}

	amIndex := strings.Index(tweet, "AM")
	if amIndex != -1 {
		briefingTime := tweet[amIndex-5 : amIndex+2]
		fmt.Println(briefingTime)
		return briefingTime
	}

	return ""
}
