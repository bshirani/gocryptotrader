package communications

import (
	"errors"

	"gocryptotrader/communications/base"
	"gocryptotrader/communications/slack"
	"gocryptotrader/communications/smsglobal"
	"gocryptotrader/communications/smtpservice"
	"gocryptotrader/communications/telegram"
)

// Communications is the overarching type across the communications packages
type Communications struct {
	base.IComm
}

// ErrNoRelayersEnabled returns when no communication relayers are enabled
var ErrNoRelayersEnabled = errors.New("no communication relayers are enabled")

// NewComm sets up and returns a pointer to a Communications object
func NewComm(cfg *base.CommunicationsConfig) (*Communications, error) {
	if !cfg.IsAnyEnabled() {
		return nil, ErrNoRelayersEnabled
	}

	var comm Communications
	if cfg.TelegramConfig.Enabled {
		Telegram := new(telegram.Telegram)
		Telegram.Setup(cfg)
		comm.IComm = append(comm.IComm, Telegram)
	}

	if cfg.SMSGlobalConfig.Enabled {
		SMSGlobal := new(smsglobal.SMSGlobal)
		SMSGlobal.Setup(cfg)
		comm.IComm = append(comm.IComm, SMSGlobal)
	}

	if cfg.SMTPConfig.Enabled {
		SMTP := new(smtpservice.SMTPservice)
		SMTP.Setup(cfg)
		comm.IComm = append(comm.IComm, SMTP)
	}

	if cfg.SlackConfig.Enabled {
		Slack := new(slack.Slack)
		Slack.Setup(cfg)
		comm.IComm = append(comm.IComm, Slack)
	}

	comm.Setup()
	return &comm, nil
}
