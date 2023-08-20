package model

import (
	ewrap "autoposting/pkg/err-wrapper"
)

type SocialNetworkName string

const (
	VK  SocialNetworkName = "VK"
	OK  SocialNetworkName = "OK"
	FB  SocialNetworkName = "FB"
	TWI SocialNetworkName = "TWI"
)

func (sn SocialNetworkName) Validate() error {
	switch sn {
	case VK, OK, FB, TWI:
		return nil
	default:
		return ewrap.Errorf("social network %s is not valid", sn)
	}
}
