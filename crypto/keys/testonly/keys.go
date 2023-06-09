// Copyright 2017 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package testonly contains code and data that should only be used by tests.
// Production code MUST NOT depend on anything in this package. This will be
// enforced by tools where possible.
package testonly

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ed25519"
)

// MustMarshalPublicPEMToDER reads a PEM-encoded public key and returns it in DER encoding.
// If an error occurs, it panics.
func MustMarshalPublicPEMToDER(keyPEM string) []byte {
	block, _ := pem.Decode([]byte(keyPEM))
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	keyDER, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		panic(err)
	}
	return keyDER
}

// SignAndVerify exercises a signer by using it to generate a signature, and
// then verifies that this signature is correct.
func SignAndVerify(signer crypto.Signer, pubKey crypto.PublicKey) error {
	hasher := crypto.SHA256
	msg := []byte("test")
	digest := sha256.Sum256(msg)

	var signature []byte
	var err error
	switch pubKey.(type) {
	case ed25519.PublicKey:
		// Ed25519 performs two passes over the data and so takes the whole message not just the digest.
		signature, err = signer.Sign(rand.Reader, msg, crypto.Hash(0))
	default:
		signature, err = signer.Sign(rand.Reader, digest[:], hasher)
	}
	if err != nil {
		return err
	}

	switch pubKey := pubKey.(type) {
	case *ecdsa.PublicKey:
		return verifyECDSA(pubKey, digest[:], signature)
	case *rsa.PublicKey:
		return verifyRSA(pubKey, digest[:], signature, hasher, hasher)
	case ed25519.PublicKey:
		return verifyEd25519(pubKey, msg, signature)
	default:
		return fmt.Errorf("unknown public key type: %T", pubKey)
	}
}

func verifyECDSA(pubKey *ecdsa.PublicKey, digest, sig []byte) error {
	var ecdsaSig struct {
		R, S *big.Int
	}

	rest, err := asn1.Unmarshal(sig, &ecdsaSig)
	if err != nil {
		return err
	}
	if len(rest) != 0 {
		return fmt.Errorf("ECDSA signature %v bytes longer than expected", len(rest))
	}

	if !ecdsa.Verify(pubKey, digest, ecdsaSig.R, ecdsaSig.S) {
		return errors.New("ECDSA signature failed verification")
	}
	return nil
}

func verifyRSA(pubKey *rsa.PublicKey, digest, sig []byte, hasher crypto.Hash, opts crypto.SignerOpts) error {
	if pssOpts, ok := opts.(*rsa.PSSOptions); ok {
		return rsa.VerifyPSS(pubKey, hasher, digest, sig, pssOpts)
	}
	return rsa.VerifyPKCS1v15(pubKey, hasher, digest, sig)
}

func verifyEd25519(pubKey ed25519.PublicKey, digest, sig []byte) error {
	if !ed25519.Verify(pubKey, digest, sig) {
		return errors.New("ed25519 signature failed verification")
	}
	return nil
}
