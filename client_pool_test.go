package smpp

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientPool(t *testing.T) {
	t.Parallel()

	require.NotNil(t, testSameServerPool)
	require.Equal(t, testPoolSize, testSameServerPool.Size())

	n := 5
	submitReportsChannel := make(chan *SubmitReport, n)
	smsList := make([]*SMS, n)
	ctx := context.Background()

	for i := 0; i < n; i++ {
		sms := &SMS{
			ID:       fmt.Sprintf("Pool-%d", i),
			Sender:   fmt.Sprintf("Pool-Trust-IQ.%d", i),
			Receiver: fmt.Sprintf("84977000%03d", i),
			Content:  RandomString(70, KeyboardCharacters),
		}
		smsList[i] = sms

		go func() {
			report := testSameServerPool.Send(ctx, sms)
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
	for _, sms := range smsList {
		sender := sms.Receiver
		content := sms.Content

		inboundSMS := testEventHandler.FindInboundSMS(sender, content)
		require.NotNil(t, inboundSMS)
	}
}

func TestClientPoolWithDifferentServers(t *testing.T) {
	t.Parallel()

	require.NotNil(t, testDifferentServersPool)
	require.Equal(t, len(testConfigs)-1, testDifferentServersPool.Size())

	n := 5
	submitReportsChannel := make(chan *SubmitReport, n)
	smsList := make([]*SMS, n)
	ctx := context.Background()

	for i := 0; i < n; i++ {
		sms := &SMS{
			ID:       fmt.Sprintf("N-%d", i),
			Sender:   fmt.Sprintf("N-Trust-IQ.%d", i),
			Receiver: fmt.Sprintf("84967000%03d", i),
			Content:  RandomString(70, KeyboardCharacters),
		}
		smsList[i] = sms

		go func() {
			report := testDifferentServersPool.Send(ctx, sms)
			submitReportsChannel <- report
		}()
	}

	// wait for delivery report and inbound sms
	WaitSeconds(20)

	for i := 0; i < n; i++ {
		submitReport := <-submitReportsChannel
		require.Nil(t, submitReport.Error)
		require.NotEmpty(t, submitReport.RemoteID)

		deliveryReport := testDifferentServersHandler.FindDeliveryReport(submitReport.RemoteID)
		require.NotNil(t, deliveryReport)
		require.True(t, deliveryReport.Delivered)
	}

	// find loop back sms
	for _, sms := range smsList {
		sender := sms.Receiver
		content := sms.Content

		inboundSMS := testDifferentServersHandler.FindInboundSMS(sender, content)
		require.NotNil(t, inboundSMS)
	}
}

func TestGetAnotherClient(t *testing.T) {
	t.Parallel()

	currentClient := testSameServerPool.clients[0]
	anotherClient := testSameServerPool.GetAnotherClient(currentClient)

	require.NotNil(t, anotherClient)
	require.NotEqual(t, currentClient, anotherClient)
}
