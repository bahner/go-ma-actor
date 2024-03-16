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

var fortuneArgs = []string{"-o"}

func reply(m *msg.Message) []byte {

	if string(m.Content) == string(replyBytes()) {
		return angryBytes()
	}

	if config.PongFortuneMode() {
		return getFortuneCookie(fortuneArgs)
	}

	return replyBytes()
}

func replyBytes() []byte {
	replyMsg := viper.GetString("mode.pong.reply")
	return []byte(replyMsg)
}

func angryBytes() []byte {
	replyMsg := viper.GetString("mode.pong.reply")
	return []byte(fmt.Sprintf("I'm doing the %s here! 😤", replyMsg))
}

// Returns a fortune cookie if the pong-fortune mode is enabled, otherwise the default reply.
func getFortuneCookie(args []string) []byte {
	// Check if the fortune command is available in the PATH
	_, err := exec.LookPath("fortune")
	if err != nil {
		log.Errorf("fortune command not found: %s", err)
		return replyBytes()
	}

	// Prepare the command with any arguments passed
	cmd := exec.Command("fortune", args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	// Execute the command
	err = cmd.Run()
	if err != nil {
		log.Errorf("error running fortune command: %s", err)
		return replyBytes()
	}

	return out.Bytes()
}
