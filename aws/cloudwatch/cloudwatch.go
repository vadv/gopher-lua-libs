// Package cloudwatch implements cloudwatch client api functionality for lua.
package cloudwatch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	lua "github.com/yuin/gopher-lua"
)

type luaClW struct {
	awsClwClient *cloudwatchlogs.CloudWatchLogs
}

func checkluaClW(L *lua.LState, n int) *luaClW {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*luaClW); ok {
		return v
	}
	L.ArgError(1, "clw expected")
	return nil
}

func newLauClW(awsProfile *string, awsRegion *string) (*luaClW, error) {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if awsProfile != nil {
		opts.Profile = *awsProfile
	}
	if awsRegion != nil {
		opts.Config = aws.Config{Region: awsRegion}
	}
	sess := session.Must(session.NewSessionWithOptions(opts))
	return &luaClW{awsClwClient: cloudwatchlogs.New(sess)}, nil
}

func (cw *luaClW) events(filterParams *cloudwatchlogs.FilterLogEventsInput, events chan *cloudwatchlogs.FilteredLogEvent, done chan error) {

	processor := func(res *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		for _, event := range res.Events {
			events <- event
		}
		if lastPage {
			close(events)
		}
		return !lastPage
	}

	done <- cw.awsClwClient.FilterLogEventsPages(filterParams, processor)
}
