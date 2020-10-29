package smpp

var smsShortASCII = &SMS{
	ID:       "sms#1",
	Sender:   "Trust-IQ.1",
	Receiver: "84987000001",
	Content:  RandomString(160, KeyboardCharacters),
}

var smsLongASCII = &SMS{
	ID:       "sms#2",
	Sender:   "Trust-IQ.2",
	Receiver: "84987000002",
	Content:  RandomString(170, KeyboardCharacters),
}

var smsShortUnicode = &SMS{
	ID:       "sms#3",
	Sender:   "Trust-IQ.3",
	Receiver: "84987000003",
	Content:  RandomString(70, UnicodeLetters),
}

var smsLongUnicode = &SMS{
	ID:       "sms#4",
	Sender:   "Trust-IQ.4",
	Receiver: "84987000004",
	Content:  RandomString(80, UnicodeLetters),
}

var smsShortHindiUnicode = &SMS{
	ID:       "sms#5",
	Sender:   "Trust-IQ.5",
	Receiver: "84987000005",
	Content:  "आता मिळवा पर्सनल लोन सुलभ इ एम आई आणि आकर्षक व्याज दराने",
}

var smsLongHindiUnicode = &SMS{
	ID:       "sms#6",
	Sender:   "Trust-IQ.6",
	Receiver: "84987000006",
	Content:  "आता मिळवा पर्सनल लोन सुलभ इ एम आई आणि आकर्षक व्याज दराने. अर्ज करण्याकरिता 59333 वर CL टाइप करुन पाठवा. नियम आणि अटी लागू",
}

var smsShortVNUnicode = &SMS{
	ID:       "sms#7",
	Sender:   "Trust-IQ.7",
	Receiver: "84987000007",
	Content:  "Đây là Tiếng Việt",
}
