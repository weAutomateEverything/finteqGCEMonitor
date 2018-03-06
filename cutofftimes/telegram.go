package cutofftimes

import (
	"context"
	"github.com/kyokomi/emoji"
	"github.com/weAutomateEverything/go2hal/remoteTelegramCommands"
	"log"
	"strconv"
	"time"
)

func (s *service) registerRemoteStream() {
	for {

		request := remoteTelegramCommands.RemoteCommandRequest{Description: "Disable GCE Cutoff time check", Name: "DisableGCECutoffTimes"}
		stream, err := s.client.RegisterCommand(context.Background(), &request)
		if err != nil {
			log.Println(err)
		} else {
			s.monitorForStreamResponse(stream)
		}
		time.Sleep(30 * time.Second)
	}
}

func (s *service) monitorForStreamResponse(client remoteTelegramCommands.RemoteCommand_RegisterCommandClient) {
	for {
		in, err := client.Recv()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(in.From)
		log.Println(in.Message)
		s.disabled = true
		var t int64
		t = 60
		if in.Message != "" {
			t, err = strconv.ParseInt(in.Message, 10, 64)
			if err != nil {
				log.Println(err)
				t = 60
			}
		}
		s.disabledTill = time.Now().Add(time.Duration(t) * time.Minute)
		s.alert.SendAlert(context.TODO(), emoji.Sprintf(":zzz: GCE service cutt-off time alerts will now sleep for %v minutes", t))
	}
}

func (s *service) registerTriggerGCECheckStream() {
	for {

		request := remoteTelegramCommands.RemoteCommandRequest{Description: "Trigger GCE Cutofftime check", Name: "GCECutofftimes"}
		stream, err := s.client.RegisterCommand(context.Background(), &request)
		if err != nil {
			log.Println(err)
		} else {
			s.monitorForGCEStreamResponse(stream)
		}
		time.Sleep(30 * time.Second)
	}
}

func (s *service) monitorForGCEStreamResponse(client remoteTelegramCommands.RemoteCommand_RegisterCommandClient) {
	s.DoCheck(true)
	s.DoCheck(false)
}
