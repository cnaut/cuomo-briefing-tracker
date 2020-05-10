package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func main() {
	config := oauth1.NewConfig(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_KEY_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
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
			timeEndIndex := strings.Index(tweet.Text, "ET")
			briefingTime := tweet.Text[timeEndIndex-8 : timeEndIndex+2]
			fmt.Println(tweet.Text)
			fmt.Println(briefingTime)

			accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
			authToken := os.Getenv("TWILIO_AUTH_TOKEN")
			urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

			// Pack up the data for our message
			msgData := url.Values{}
			msgData.Set("To", "+17183442807")
			msgData.Set("From", "+12513068087")
			msgData.Set("Body", "Watch Cuomo daily briefing at "+briefingTime)
			msgDataReader := *strings.NewReader(msgData.Encode())

			// Create HTTP request client
			client := &http.Client{}
			req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
			req.SetBasicAuth(accountSid, authToken)
			req.Header.Add("Accept", "application/json")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			// Make HTTP POST request and return message SID
			resp, _ := client.Do(req)
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				var data map[string]interface{}
				decoder := json.NewDecoder(resp.Body)
				err := decoder.Decode(&data)
				if err == nil {
					fmt.Println(data["sid"])
				}
			} else {
				fmt.Println(resp.Status)
			}

			googleCredentials, err := ioutil.ReadFile("credentials.json")
			if err != nil {
				log.Fatalf("Unable to read client secret file: %v", err)
			}

			googleConfig, err := google.ConfigFromJSON(googleCredentials, calendar.CalendarReadonlyScope)
			if err != nil {
				log.Fatalf("Unable to parse client secret file to config: %v", err)
			}

			googleClient := getGoogleClient(googleConfig)

			srv, err := calendar.New(googleClient)
			if err != nil {
				log.Fatalf("Unable to retrieve Calendar client: %v", err)
			}

			timeNow := time.Now()
			briefingStartTimeFormatted := fmt.Sprintf("%d-%02d-%02dT"+briefingTime[0:5]+":00-04:00", timeNow.Year(), timeNow.Month(), timeNow.Day())
			briefingTimeParsed, _ := time.Parse(
				time.RFC3339,
				briefingStartTimeFormatted)

			briefingEndTimeFormatted := briefingTimeParsed.Add(time.Hour).Format("2006-01-02T15:04:05-07:00")
			fmt.Println(string(briefingStartTimeFormatted))
			fmt.Println(string(briefingEndTimeFormatted))

			event := &calendar.Event{
				Summary:     "Governor Cuomo Daily Briefing",
				Location:    "ny.gov",
				Description: "NY Governor Andrew Cuomo's daily coronavirus briefing ",
				Start: &calendar.EventDateTime{
					DateTime: briefingStartTimeFormatted,
					TimeZone: "America/New_York",
				},

				End: &calendar.EventDateTime{
					DateTime: briefingEndTimeFormatted,
					TimeZone: "America/New_York",
				},
			}

			calendarID := "primary"
			event, err = srv.Events.Insert(calendarID, event).Do()
			if err != nil {
				log.Fatalf("Unable to create event. %v\n", err)
			}
			fmt.Printf("Event created: %s\n", event.HtmlLink)

			return
		}
	}
}
