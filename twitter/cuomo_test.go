package twitter

import "testing"

func TestFindBriefingTimeInString(t *testing.T) {
	tweet := "I will be giving my daily #Coronavirus briefing at 11:30AM. Watch here:"
	briefingTime := FindBriefingTimeInString(tweet)
	if briefingTime != "11:30AM" {
		t.Errorf("Briefing time returned was %q, Expected 11:30AM", briefingTime)
	}

	tweet = "I will be giving my daily #Coronavirus briefing at 12:00PM. Watch here:"
	briefingTime = FindBriefingTimeInString(tweet)
	if briefingTime != "12:00PM" {
		t.Errorf("Briefing time returned was %q, Expected 12:00PM", briefingTime)
	}

	tweet = "I will be giving my daily #Coronavirus briefing at 11:30AM ET. Watch here:"
	briefingTime = FindBriefingTimeInString(tweet)
	if briefingTime != "11:30AM ET" {
		t.Errorf("Briefing time returned was %q, Expected 11:30AM ET", briefingTime)
	}

	tweet = "Do your part. Show respect. Wear a mask."
	briefingTime = FindBriefingTimeInString(tweet)
	if briefingTime != "" {
		t.Errorf("Briefing time returned was %q, Expected no briefing time", briefingTime)
	}
}
