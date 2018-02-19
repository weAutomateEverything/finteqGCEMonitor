package cutofftimes

import (
	"github.com/zamedic/go2hal/remoteTelegramCommands"
	"time"
	"context"
	"log"
	"strconv"
	"github.com/kyokomi/emoji"
)

func (s *service) registerRemoteStream() {
	for {
		name := "DisableGCECuttoffInward"
		if s.inward{
			name = "DisableGCECuttoffOutward"
		}
		request := remoteTelegramCommands.RemoteCommandRequest{Description: "Disable GCE Cutoff time checl", Name: name}
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
			t, err = strconv.ParseInt(in.Message,10,64)
			if err != nil {
				log.Println(err)
				t = 60
			}
		}
		s.disabledTill = time.Now().Add(time.Duration(t) * time.Minute)
		s.alert.SendAlert(emoji.Sprintf(":zzz: smoke tests will now sleep for %v minutes", t))
	}
}
