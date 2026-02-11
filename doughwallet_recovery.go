package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/cyclone-github/base58"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
)

/*
Program:    Dough Wallet Recovery Tool
Author:     Cyclone
Version:    0.1.3
Date:       2025/02/11
GitHub:     github.com/cyclone-github/doughwallet_recovery
License:	GPL-2 - https://github.com/cyclone-github/doughwallet_recovery?tab=GPL-2.0-1-ov-file

Dough Wallet uses a non-standard BIP32 hardened derivation flag
Standard BIP32 uses 0x80000000
Dough Wallet uses 0x9e000000
This means all hardened child keys are derived with different indices, producing completely different keys from a standard BIP32 wallet
*/

func versionFunc() {
	fmt.Fprintln(os.Stderr, "Dough Wallet Recovery v0.1.3; 2026-02-11\nhttps://github.com/cyclone-github/")
}

const bip32Prime uint32 = 0x9e000000

const bip32SeedKey = "Bitcoin seed"

// Dogecoin mainnet address parameters
const (
	dogePubKeyVersion byte = 30  // 0x1e - P2PKH addresses start with 'D'
	dogeWIFVersion    byte = 158 // 0x9e - WIF private keys start with 'Q'
)

// secp256k1 curve order N
var curveN, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

// encode payload with a version byte and 4-byte checksum
func base58CheckEncode(version byte, payload []byte) string {
	data := make([]byte, 0, 1+len(payload)+4)
	data = append(data, version)
	data = append(data, payload...)

	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	data = append(data, hash2[:4]...)

	return base58.StdEncoding.EncodeToString(data)
}

// RIPEMD-160(SHA-256(data))
func hash160(data []byte) []byte {
	sha := sha256.Sum256(data)
	r := ripemd160.New()
	r.Write(sha[:])
	return r.Sum(nil)
}

// generate compressed public key (33 bytes) from 32-byte private key
func compressedPubKey(privKeyBytes []byte) []byte {
	if len(privKeyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(privKeyBytes):], privKeyBytes)
		privKeyBytes = padded
	}
	privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)
	return privKey.PubKey().SerializeCompressed()
}

/*
ckdPriv performs BIP32 private child key derivation
This implements the standard BIP32 CKD function, but Dough Wallet uses a non-standard hardened flag (0x9e000000 instead of 0x80000000)
The flag value is encoded in the index parameter
For hardened derivation (index & 0x9e000000 != 0): HMAC-SHA512(Key=chainCode, Data=0x00 || privKey || index)
For normal derivation: HMAC-SHA512(Key=chainCode, Data=compressedPubKey || index)
*/

// return derived child private key and chain code
func ckdPriv(key, chainCode []byte, index uint32) ([]byte, []byte) {
	var data []byte

	if index&bip32Prime != 0 {
		// hardened derivation: 0x00 || key || index (big-endian)
		data = make([]byte, 0, 1+32+4)
		data = append(data, 0x00)
		data = append(data, key...)
	} else {
		// normal derivation: compressed_pubkey || index (big-endian)
		pubKey := compressedPubKey(key)
		data = make([]byte, 0, 33+4)
		data = append(data, pubKey...)
	}

	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	data = append(data, indexBytes...)

	mac := hmac.New(sha512.New, chainCode)
	mac.Write(data)
	I := mac.Sum(nil)

	// child key = (parse256(IL) + kpar) mod n
	IL := new(big.Int).SetBytes(I[:32])
	kPar := new(big.Int).SetBytes(key)
	IL.Add(IL, kPar)
	IL.Mod(IL, curveN)

	// left-pad the result to 32 bytes
	newKey := make([]byte, 32)
	b := IL.Bytes()
	copy(newKey[32-len(b):], b)

	// new chain code is the right 32 bytes of I
	newChainCode := make([]byte, 32)
	copy(newChainCode, I[32:])

	return newKey, newChainCode
}

// convert compressed public key to Dogecoin P2PKH address
func pubKeyToAddress(pubKey []byte) string {
	h := hash160(pubKey)
	return base58CheckEncode(dogePubKeyVersion, h)
}

// convert 32-byte private key to Wallet Import Format (compressed)
func privKeyToWIF(privKey []byte) string {
	payload := make([]byte, 0, 33)
	payload = append(payload, privKey...)
	payload = append(payload, 0x01)
	return base58CheckEncode(dogeWIFVersion, payload)
}

func main() {
	version := flag.Bool("version", false, "Program version")
	cyclone := flag.Bool("cyclone", false, "Dough Wallet Recovery")
	count := flag.Uint("count", 1, "Number of addresses to generate)")
	flag.Parse()

	// run sanity checks for special flags
	if *version {
		versionFunc()
		os.Exit(0)
	}
	if *cyclone {
		decodedStr, err := base64.StdEncoding.DecodeString("Q29kZWQgYnkgY3ljbG9uZSA7KQo=")
		if err != nil {
			fmt.Fprintln(os.Stderr, "--> Cannot decode base64 string. <--")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, string(decodedStr))
		os.Exit(0)
	}

	fmt.Println(" --------------------------- ")
	fmt.Println("|   Dough Wallet Recovery   |")
	fmt.Println("| github.com/cyclone-github |")
	fmt.Println(" --------------------------- ")
	fmt.Println()
	fmt.Print("Enter your Dough Wallet's 12-word recovery phrase:\n")

	reader := bufio.NewReader(os.Stdin)
	phrase, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	phrase = strings.TrimSpace(phrase)
	phrase = strings.ToLower(phrase)

	// split into words and validate count
	words := strings.Fields(phrase)
	if len(words) != 12 {
		fmt.Fprintf(os.Stderr, "Error: expected 12 words, got %d\n", len(words))
		os.Exit(1)
	}

	phrase = strings.Join(words, " ")

	// validate each word against BIP39 en wordlist
	wordList := bip39.GetWordList()
	wordSet := make(map[string]bool, len(wordList))
	for _, w := range wordList {
		wordSet[w] = true
	}

	for i, w := range words {
		if !wordSet[w] {
			fmt.Fprintf(os.Stderr, "Error: word #%d \"%s\" is not in the BIP39 wordlist\n", i+1, w)
			os.Exit(1)
		}
	}

	// generate a 64-byte seed from the mnemonic
	// uses PBKDF2(mnemonic, "mnemonic", 2048, 64, HMAC-SHA512)
	seed := bip39.NewSeed(phrase, "")

	// BIP32: master private key and chain code
	// master = HMAC-SHA512(Key="Bitcoin seed", Data=seed)
	mac := hmac.New(sha512.New, []byte(bip32SeedKey))
	mac.Write(seed)
	I := mac.Sum(nil)

	masterKey := make([]byte, 32)
	masterChainCode := make([]byte, 32)
	copy(masterKey, I[:32])
	copy(masterChainCode, I[32:])

	// account key: m/0' (hardened with non-standard flag 0x9e000000)
	acctKey, acctChainCode := ckdPriv(masterKey, masterChainCode, 0|bip32Prime)

	fmt.Println()
	// external (Receive) Chain: m/0'/0/n
	extChainKey, extChainCode := ckdPriv(acctKey, acctChainCode, 0)

	fmt.Println("External (Receive) Chain: m/0'/0/n")
	fmt.Println()
	for n := uint32(0); n < uint32(*count); n++ {
		childKey, _ := ckdPriv(extChainKey, extChainCode, n)
		childPub := compressedPubKey(childKey)
		fmt.Printf("  m/0'/0/%d\n", n)
		fmt.Printf("    Address:     %s\n", pubKeyToAddress(childPub))
		fmt.Printf("    Private Key: %s\n", privKeyToWIF(childKey))
		//fmt.Println()
	}

	// internal (Change) Chain: m/0'/1/n
	intChainKey, intChainCode := ckdPriv(acctKey, acctChainCode, 1)

	fmt.Println()
	fmt.Println("Internal (Change) Chain: m/0'/1/n")
	fmt.Println()
	for n := uint32(0); n < uint32(*count); n++ {
		childKey, _ := ckdPriv(intChainKey, intChainCode, n)
		childPub := compressedPubKey(childKey)
		fmt.Printf("  m/0'/1/%d\n", n)
		fmt.Printf("    Address:     %s\n", pubKeyToAddress(childPub))
		fmt.Printf("    Private Key: %s\n", privKeyToWIF(childKey))
		//fmt.Println()
	}
	fmt.Println()
}
