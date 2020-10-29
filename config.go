package smpp

import (
	"context"
	"strings"
	"time"

	"github.com/tsocial/env"
	"github.com/tsocial/go-smpp/smpp/pdu"
	"github.com/tsocial/logger"
)

// Config contains configuration for SMPP client
type Config struct {
	ID               int
	Address          string
	Username         string
	Password         string
	SystemType       string
	SMSPerSeconds    int
	ProtocolID       uint8
	SourceTON        uint8
	SourceNPI        uint8
	DestTON          uint8
	DestNPI          uint8
	ReplaceIfPresent uint8
	ResponseTimeout  time.Duration
	AddressRange     *pdu.AddressRange
	Debug            bool
}

const (
	defaultAddresses = "localhost:2775,localhost:2777,localhost:2779"
	defaultUsernames = "smpp,smpp,smpp"
	defaultPasswords = "password,password,password"

	defaultSystemType             = "socdel1"
	defaultSMSPerSeconds          = 100
	defaultProtocolID             = 0
	defaultSourceTON              = 0
	defaultSourceNPI              = 0
	defaultDestTON                = 0
	defaultDestNPI                = 0
	defaultReplaceIfPresent       = 0
	defaultResponseTimeoutSeconds = 1
	defaultAddrTON                = 0
	defaultAddrNPI                = 0
	defaultAddrRange              = ""
	delimiter                     = ","
)

// LoadConfigsFromEnvironment loads configuration from environment variables
func LoadConfigsFromEnvironment() []*Config {
	addresses := strings.Split(env.EVString("SMPP_ADDRESSES", defaultAddresses), delimiter)
	usernames := strings.Split(env.EVString("SMPP_USERNAMES", defaultUsernames), delimiter)
	passwords := strings.Split(env.EVString("SMPP_PASSWORDS", defaultPasswords), delimiter)
	configs := make([]*Config, len(addresses))
	responseTimeoutSeconds := env.EVInt("SMPP_RESPONSE_TIMEOUT_SECONDS", defaultResponseTimeoutSeconds)
	for idx, address := range addresses {
		configs[idx] = &Config{
			ID:               idx,
			Address:          address,
			Username:         usernames[idx],
			Password:         passwords[idx],
			SystemType:       env.EVString("SMPP_SYSTEM_TYPE", defaultSystemType),
			SMSPerSeconds:    env.EVInt("SMPP_SMS_PER_SECONDS", defaultSMSPerSeconds),
			ProtocolID:       uint8(env.EVInt("SMPP_PROTOCOL_ID", defaultProtocolID)),
			SourceTON:        uint8(env.EVInt("SMPP_SOURCE_TON", defaultSourceTON)),
			SourceNPI:        uint8(env.EVInt("SMPP_SOURCE_NPI", defaultSourceNPI)),
			DestTON:          uint8(env.EVInt("SMPP_DEST_TON", defaultDestTON)),
			DestNPI:          uint8(env.EVInt("SMPP_DEST_NPI", defaultDestNPI)),
			ReplaceIfPresent: uint8(env.EVInt("SMPP_REPLACE_IF_PRESENT", defaultReplaceIfPresent)),
			ResponseTimeout:  time.Duration(responseTimeoutSeconds) * time.Second,
			AddressRange:     loadAddressRange(),
			Debug:            env.EVBool("DEBUG", true),
		}
		logger.Println(context.Background(), "smpp config loaded", configs[idx])
	}

	return configs
}

func loadAddressRange() *pdu.AddressRange {
	if !env.HasEnvironmentKey("SMPP_ADDR_RANGE") {
		return nil
	}

	address := env.EVString("SMPP_ADDR_RANGE", defaultAddrRange)
	ton := uint8(env.EVInt("SMPP_ADDR_TON", defaultAddrTON))
	npi := uint8(env.EVInt("SMPP_ADDR_NPI", defaultAddrNPI))

	return &pdu.AddressRange{
		TON:     ton,
		NPI:     npi,
		Address: address,
	}
}
