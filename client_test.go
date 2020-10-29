package smpp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSMPPClient(t *testing.T) {
	t.Parallel()

	require.NotNil(t, testClient)

	shortSMSList := []*SMS{smsShortASCII, smsShortUnicode, smsShortHindiUnicode, smsShortVNUnicode}
	longSMSList := []*SMS{smsLongASCII, smsLongUnicode, smsLongHindiUnicode}
	smsList := append(shortSMSList, longSMSList...)
	n := len(smsList)

	ctx := context.Background()
	submitReportsChannel := make(chan *SubmitReport, n)
	for i := 0; i < n; i++ {
		sms := smsList[i]
		go func() {
			report := testClient.Send(ctx, sms)
			submitReportsChannel <- report
		}()
	}

	// wait for delivery report and inbound sms
	WaitSeconds(20)

	for i := 0; i < n; i++ {
		submitReport := <-submitReportsChannel
		require.Nil(t, submitReport.Error)
		require.NotEmpty(t, submitReport.RemoteID)

		deliveryReport := testEventHandler.FindDeliveryReport(submitReport.RemoteID)
		require.NotNil(t, deliveryReport)
		require.True(t, deliveryReport.Delivered)
	}

	// find loop back sms
	for _, sms := range shortSMSList {
		sender := sms.Receiver
		content := sms.Content

		inboundSMS := testEventHandler.FindInboundSMS(sender, content)
		require.NotNil(t, inboundSMS)
	}
}
