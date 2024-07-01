package fakes

import (
	"encoding/json"

	"github.com/juanfcgarcia/gostori/internal/notification"
)

type DummyNotifier struct {
	Receiver DummyNotifierReceiver
}

type DummyNotifierReceiver struct {
	MsgList []string
}

func NewDummyNotifier() *DummyNotifier {
	return &DummyNotifier{}
}

func (dn *DummyNotifier) SendSummaryEmail(params notification.SendSummaryEmailParams) error {
	paramJson, _ := json.Marshal(params)
	dn.Receiver.MsgList = append(dn.Receiver.MsgList, string(paramJson))
	return nil
}
