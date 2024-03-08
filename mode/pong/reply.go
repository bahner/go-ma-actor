package pong

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	replyBytes  = []byte(viper.GetString("mode.pong.reply"))
	angryBytes  = []byte(fmt.Sprintf("I'm doing the %s here! 😤", replyBytes))
	fortuneArgs = []string{"-s"}
)

func reply(m *msg.Message) []byte {
	if string(m.Content) == string(replyBytes) {
		return angryBytes
	}

	if config.PongFortuneMode() {
		return getFortuneCookie(fortuneArgs)
	}

	return replyBytes
}

// Returns a fortune cookie if the pong-fortune mode is enabled, otherwise the default reply.
func getFortuneCookie(args []string) []byte {
	// Check if the fortune command is available in the PATH
	_, err := exec.LookPath("fortune")
	if err != nil {
		log.Errorf("fortune command not found: %s", err)
		return replyBytes
	}

	// Prepare the command with any arguments passed
	cmd := exec.Command("fortune", args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	// Execute the command
	err = cmd.Run()
	if err != nil {
		log.Errorf("error running fortune command: %s", err)
		return replyBytes
	}

	return out.Bytes()
}
