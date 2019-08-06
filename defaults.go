package slacker

import (
	"context"
	"errors"
	"github.com/nlopes/slack"
)

// ClientOption an option for client values
type ClientOption func(*ClientDefaults)

// WithDebug sets debug toggle
func WithDebug(debug bool) ClientOption {
	return func(defaults *ClientDefaults) {
		defaults.Debug = debug
	}
}

// ClientDefaults configuration
type ClientDefaults struct {
	Debug bool
}

func newClientDefaults(options ...ClientOption) *ClientDefaults {
	config := &ClientDefaults{
		Debug: false,
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// ReplyOption an option for reply values
type ReplyOption func(*ReplyDefaults)

// WithAttachments sets message attachments
func WithAttachments(attachments []slack.Attachment) ReplyOption {
	return func(defaults *ReplyDefaults) {
		defaults.Attachments = attachments
	}
}

// WithBlocks sets message blocks
func WithBlocks(blocks []slack.Block) ReplyOption {
	return func(defaults *ReplyDefaults) {
		defaults.Blocks = blocks
	}
}

// ReplyDefaults configuration
type ReplyDefaults struct {
	Attachments []slack.Attachment
	Blocks      []slack.Block
}

func newReplyDefaults(options ...ReplyOption) *ReplyDefaults {
	config := &ReplyDefaults{
		Attachments: []slack.Attachment{},
		Blocks:      []slack.Block{},
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// DefaultEventHandler it the default event handler.
var DefaultEventHandler = func(ctx context.Context, s *Slacker, msg slack.RTMEvent) error {
	switch event := msg.Data.(type) {
	case *slack.ConnectedEvent:
		if s.initHandler == nil {
			return nil
		}
		go s.initHandler()

	case *slack.MessageEvent:
		if s.isFromBot(event) {
			return nil
		}

		if !s.isBotMentioned(event) && !s.isDirectMessage(event) {
			return nil
		}
		go s.handleMessage(ctx, event)

	case *slack.RTMError:
		if s.errorHandler == nil {
			return nil
		}
		go s.errorHandler(event.Error())

	case *slack.InvalidAuthEvent:
		return errors.New(invalidToken)

	default:
		if s.fallbackEventHandler == nil {
			return nil
		}
		go s.fallbackEventHandler(event)
	}

	return nil
}