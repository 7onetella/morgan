package ssha

// MIT License

// Copyright (c) 2019 7onetella

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

/*
   All credit goes to https://medium.com/@ferdinand.neman/ssha-password-hash-with-golang-7d79d792bd3d
*/

// Encode encodes the []byte of raw password
func Encode(rawPassPhrase []byte) ([]byte, error) {
	hash := makeSSHAHash(rawPassPhrase, makeSalt())
	b64 := base64.StdEncoding.EncodeToString(hash)
	return []byte(fmt.Sprintf("{SSHA}%s", b64)), nil
}

// Matches matches the encoded password and the raw password
func Matches(encodedPassPhrase, rawPassPhrase []byte) bool {
	//strip the {SSHA}
	eppS := string(encodedPassPhrase)[6:]
	hash, err := base64.StdEncoding.DecodeString(eppS)
	if err != nil {
		return false
	}
	salt := hash[len(hash)-4:]

	sha := sha1.New()
	sha.Write(rawPassPhrase)
	sha.Write(salt)
	sum := sha.Sum(nil)

	if bytes.Compare(sum, hash[:len(hash)-4]) != 0 {
		return false
	}
	return true
}

// makeSalt make a 4 byte array containing random bytes.
func makeSalt() []byte {
	sbytes := make([]byte, 4)
	rand.Read(sbytes)
	return sbytes
}

// makeSSHAHash make hasing using SHA-1 with salt. This is not the final output though. You need to append {SSHA} string with base64 of this hash.
func makeSSHAHash(passphrase, salt []byte) []byte {
	sha := sha1.New()
	sha.Write(passphrase)
	sha.Write(salt)

	h := sha.Sum(nil)
	return append(h, salt...)
}
