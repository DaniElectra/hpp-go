package hpp

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// GeneratePasswordSignature generates a password signature for the given buffer
func GeneratePasswordSignature(pid uint32, password string, buffer []byte) []byte {
	passwordBytes := []byte(password)

	passwordSignatureKey := DeriveKerberosKey(pid, passwordBytes)

	calculatedPasswordSignature := calculateSignature(buffer, passwordSignatureKey)

	return calculatedPasswordSignature
}

// GenerateAccessKeySignature genarates an access key signature for the given buffer
func GenerateAccessKeySignature(accessKey string, buffer []byte) ([]byte, error) {
	accessKeyBytes, err := hex.DecodeString(accessKey)
	if err != nil {
		return nil, errors.New("Failed to decode access key from server: " + err.Error())
	}

	calculatedAccessKeySignature := calculateSignature(buffer, accessKeyBytes)

	return calculatedAccessKeySignature, nil
}

func calculateSignature(buffer []byte, key []byte) []byte {
	mac := hmac.New(md5.New, key)
	mac.Write(buffer)
	hmac := mac.Sum(nil)

	return hmac
}

// DeriveKerberosKey derives a users kerberos encryption key based on their PID and password
func DeriveKerberosKey(pid uint32, password []byte) []byte {
	for i := 0; i < 65000+int(pid)%1024; i++ {
		password = MD5Hash(password)
	}

	return password
}
