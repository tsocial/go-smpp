package smpp

import (
	"context"
	"sync"

	"github.com/tsocial/logger"
)

// IEventHandler contains functions to handle events from SMPP server
type IEventHandler interface {
	HandleDeliveryReport(ctx context.Context, report *DeliveryReport)
	HandleInboundSMS(ctx context.Context, sms *SMS)
}

// SampleEventHandler is a sample event handler that can be used for test
type SampleEventHandler struct {
	deliveryReports []*DeliveryReport
	inboundSMS      []*SMS
	sync.RWMutex
}

// NewSampleEventHandler creates a new sample event handler
func NewSampleEventHandler() *SampleEventHandler {
	return &SampleEventHandler{}
}

// HandleDeliveryReport handles delivery report for outbound sms
func (handler *SampleEventHandler) HandleDeliveryReport(ctx context.Context, report *DeliveryReport) {
	handler.Lock()
	defer handler.Unlock()

	logger.Println(ctx, "received delivery report:", report)
	handler.deliveryReports = append(handler.deliveryReports, report)
}

// HandleInboundSMS handles inbound sms
func (handler *SampleEventHandler) HandleInboundSMS(ctx context.Context, sms *SMS) {
	handler.Lock()
	defer handler.Unlock()

	logger.Println(ctx, "received inbound sms:", sms)
	handler.inboundSMS = append(handler.inboundSMS, sms)
}

// FindDeliveryReport finds delivery report by remote id
func (handler *SampleEventHandler) FindDeliveryReport(remoteID string) *DeliveryReport {
	handler.RLock()
	defer handler.RUnlock()

	for _, report := range handler.deliveryReports {
		if report.RemoteID == remoteID {
			return report
		}
	}
	return nil
}

// FindInboundSMS finds inbound sms by sender and content
func (handler *SampleEventHandler) FindInboundSMS(sender string, content string) *SMS {
	handler.RLock()
	defer handler.RUnlock()

	for _, sms := range handler.inboundSMS {
		if sms.Sender == sender && sms.Content == content {
			return sms
		}
	}
	return nil
}
