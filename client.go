package smpp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tsocial/go-smpp/smpp"
	"github.com/tsocial/go-smpp/smpp/pdu"
	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
	"github.com/tsocial/go-smpp/smpp/pdu/pdutext"
	"github.com/tsocial/logger"
	"golang.org/x/time/rate"
)

// Client is a SMPP client to send SMS
type Client struct {
	config      *Config
	handler     IEventHandler
	transceiver *smpp.Transceiver
}

// NewClient creates a new client
func NewClient(ctx context.Context, config *Config, eventHandler IEventHandler) (*Client, error) {
	client := &Client{
		config:  config,
		handler: eventHandler,
	}

	client.transceiver = &smpp.Transceiver{
		Addr:         config.Address,
		User:         config.Username,
		Passwd:       config.Password,
		SystemType:   config.SystemType,
		AddressRange: config.AddressRange,
		Handler:      client.Receive,
		RespTimeout:  config.ResponseTimeout,
	}

	if config.SMSPerSeconds > 0 {
		client.transceiver.RateLimiter = rate.NewLimiter(rate.Limit(config.SMSPerSeconds), 1)
	}

	csChannel := client.transceiver.Bind()

	cs := <-csChannel
	logger.PrintInfo(ctx, fmt.Sprintf("smpp connection status: %s", cs.Status()))

	err := cs.Error()
	if err != nil {
		return nil, err
	}

	go func() {
		for cs := range csChannel {
			logger.PrintInfo(ctx, fmt.Sprintf("smpp connection status: %s", cs.Status()))
		}
	}()

	return client, nil
}

// Disconnect closes connection to SMPP server
func (client *Client) Disconnect(ctx context.Context) {
	err := client.transceiver.Close()
	logger.PrintError(ctx, err)
}

// Send sends outbound sms
func (client *Client) Send(ctx context.Context, sms *SMS) *SubmitReport {
	config := client.config
	message := &smpp.ShortMessage{
		Src:                  sms.Sender,
		Dst:                  sms.Receiver,
		Register:             pdufield.FinalDeliveryReceipt,
		ProtocolID:           config.ProtocolID,
		SourceAddrNPI:        config.SourceNPI,
		SourceAddrTON:        config.SourceTON,
		DestAddrNPI:          config.DestNPI,
		DestAddrTON:          config.DestTON,
		ReplaceIfPresentFlag: config.ReplaceIfPresent,
	}

	isUnicode := ContainsUnicodeChar(sms.Content)
	if isUnicode {
		message.Text = pdutext.UCS2(sms.Content)
	} else {
		message.Text = pdutext.Raw(sms.Content)
	}

	numberChars := RuneLength(sms.Content)

	submit := client.transceiver.SubmitLongMsg
	if (isUnicode && numberChars <= 70) || (!isUnicode && numberChars <= 160) {
		submit = client.transceiver.Submit
	}

	if client.config.Debug {
		logger.PrintInfo(ctx, "sending sms", isUnicode, numberChars, message)
	}

	startTime := time.Now()
	message, err := submit(message)
	report := &SubmitReport{
		ServerID: client.config.ID,
		SMSID:    sms.ID,
		Error:    err,
	}

	if err == nil {
		report.RemoteID = message.RespID()
	}

	if client.config.Debug {
		runTime := time.Since(startTime)
		logger.PrintInfo(ctx, "sms_time_elapsed", runTime.Nanoseconds(), runTime)
	}
	return report
}

// Receive receives inbound sms or delivery status
func (client *Client) Receive(payload pdu.Body) {
	go client.handlePDU(context.Background(), payload)
}

// handlePDU handle PDU receive from server
func (client *Client) handlePDU(ctx context.Context, payload pdu.Body) {
	header := payload.Header()

	if client.config.Debug {
		logger.PrintInfo(ctx, "receive pdu", header, header.ID.Group(), header.Key(), payload)
	}

	if header.ID == pdu.DeliverSMID {
		fields := payload.Fields()

		esmClass, err := strconv.Atoi(fields[pdufield.ESMClass].String())
		if err != nil {
			logger.PrintError(ctx, "cannot parse esm class:", err)
			return
		}

		if esmClass == 4 {
			report, err := client.parseDeliveryReport(fields)
			if err != nil {
				logger.PrintError(ctx, "cannot parse delivery report:", err)
				return
			}

			client.handler.HandleDeliveryReport(ctx, report)
			return
		}

		var dataCoding pdutext.DataCoding = 0x00
		rawDataCoding := fields[pdufield.DataCoding].Bytes()
		if len(rawDataCoding) == 1 {
			dataCoding = pdutext.DataCoding(rawDataCoding[0])
		}
		rawData := fields[pdufield.ShortMessage].Bytes()
		contentBytes := pdutext.Decode(dataCoding, rawData)
		sms := &SMS{
			Sender:   fields[pdufield.SourceAddr].String(),
			Receiver: fields[pdufield.DestinationAddr].String(),
			Content:  string(contentBytes),
		}
		client.handler.HandleInboundSMS(ctx, sms)
	} else {
		logger.PrintInfo(ctx, "receive unexpected pdu", header, payload)
	}
}

// parseDeliveryReport parses delivery report from pdufield
// http://www.smssolutions.net/tutorials/smpp/smppdeliveryreports/
// id:IIIIII sub:SSS dlvrd:DDD submit date:YYMMDDhhmm done date:YYMMDDhhmm stat:DDDDDDD err:E Text: ...
// id:c449ab9744f47b6af1879e49e75e4f40 sub:001 dlvrd:0 submit date:0610191018 done date:0610191018 stat:ACCEPTD err:0 text:sms content
func (client *Client) parseDeliveryReport(fields pdufield.Map) (*DeliveryReport, error) {
	message := fields[pdufield.ShortMessage]
	if message == nil {
		return nil, errors.New("short message is nil")
	}

	msgParts := strings.Split(message.String(), " ")
	remoteID := string([]rune(msgParts[0])[3:])
	delivered, err := strconv.Atoi(string([]rune(msgParts[2])[6:]))
	if err != nil {
		return nil, err
	}

	return &DeliveryReport{client.config.ID, remoteID, delivered == 1}, nil
}
