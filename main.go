package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nlopes/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack/slackevents"

	qiita "./lib"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// verify token
	VERIFY_TOKEN := os.Getenv("SLACK_BOT_VERIFY_TOKEN")
	reqBody := request.Body
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(reqBody),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{
				VerificationToken: VERIFY_TOKEN,
			},
		),
	)
	if err != nil {
		fmt.Print(err)
		return events.APIGatewayProxyResponse{}, err
	}

	// verify challenge param
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(reqBody), &r)
		if err != nil {
			log.Print(err)
			return events.APIGatewayProxyResponse{}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       r.Challenge,
		}, nil
	}

	// handling events
	OAUTH_TOKEN := os.Getenv("SLACK_BOT_USER_ACCESS_TOKEN")
	api := slack.New(OAUTH_TOKEN)
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			handleAppMentionEvent(api, ev)
		case *slackevents.MessageEvent:
			if ev.Channel == "im" {
				reply := "DM"
				api.PostMessage(ev.Channel, slack.MsgOptionText(reply, false))
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func handleAppMentionEvent(api *slack.Client, ev *slackevents.AppMentionEvent) {
	reply := getRecentUserItems()
	api.PostMessage(ev.Channel, slack.MsgOptionText(reply, false))
}

func getRecentUserItems() string {
	qii, err := qiita.New("https://qiita.com/api/v2", "0d9e28484a401c818f7170af5021e0af40b7a780", nil)
	if err != nil {
		panic(err)
	}

	items, err := qii.GetUserItems(context.Background(), "lelouch99v", 1, 100)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(len(items))
}
