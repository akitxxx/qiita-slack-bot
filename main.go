package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack/slackevents"
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
	reply := ev.Text
	api.PostMessage(ev.Channel, slack.MsgOptionText(reply, false))
}
