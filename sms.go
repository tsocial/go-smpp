package smpp

import "fmt"

// SMS contains data of all sms type
type SMS struct {
	ID       string
	Sender   string
	Receiver string
	Content  string
}

// SubmitReport contains outbound sms submit report data
type SubmitReport struct {
	ServerID int
	SMSID    string
	RemoteID string
	Error    error
}

// DeliveryReport contains outbound sms delivery report data
type DeliveryReport struct {
	ServerID  int
	RemoteID  string
	Delivered bool
}

func (sms *SMS) String() string {
	return fmt.Sprintf("[id: %s, sender: %s, receiver: %s, content: %s",
		sms.ID, sms.Sender, sms.Receiver, sms.Content)
}

func (report *SubmitReport) String() string {
	return fmt.Sprintf("[sms_id: %s, remote_id: %s, error: %s]",
		report.SMSID, report.RemoteID, report.Error)
}

func (report *DeliveryReport) String() string {
	return fmt.Sprintf("[remote_id: %s, delivered: %t]",
		report.RemoteID, report.Delivered)
}
