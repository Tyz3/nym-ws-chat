package nym_util

import (
	"github.com/btcsuite/btcd/btcutil/base58"
	"strings"
)

func NymAddressToBytes(address string) [96]byte {
	args := strings.Split(address, "@")
	clientArgs := strings.Split(args[0], ".")

	gatewayIdentityKey := base58.Decode(args[1])
	userIdentityKey := base58.Decode(clientArgs[0])
	userEncryptionKey := base58.Decode(clientArgs[1])

	var res [96]byte

	copy(res[0:32], userIdentityKey)
	copy(res[32:64], userEncryptionKey)
	copy(res[64:96], gatewayIdentityKey)

	return res
}

func NymAddressFromBytes(address []byte) string {
	var sb strings.Builder

	userIdentityKey := base58.Encode(address[0:32])
	userEncryptionKey := base58.Encode(address[32:64])
	gatewayIdentityKey := base58.Encode(address[64:96])

	sb.WriteString(userIdentityKey)
	sb.WriteString(".")
	sb.WriteString(userEncryptionKey)
	sb.WriteString("@")
	sb.WriteString(gatewayIdentityKey)

	return sb.String()
}
