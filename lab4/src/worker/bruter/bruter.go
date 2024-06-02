package bruter

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

type Algorithm int
var wordlist [][]byte

const (
    SHA256 Algorithm = iota 
)

type BruteforceResult struct {
    Input string
    Recovered []byte
    Err error
}

type Bruter struct {
    algorithm Algorithm
    hash []byte 
    recovered string
    wordlist [][]byte
}

func (br *Bruter)Brute() ([]byte, error) {
    var hashFunc func([]byte) []byte

    switch br.algorithm {
    case SHA256:
        hashFunc = func(msg []byte)[]byte {
            res := sha256.Sum256(msg)
            return res[:]
        }
    }
    for _, word := range wordlist {
        if bytes.Equal(hashFunc(word), br.hash) {
            return word, nil
        }
    }
    return nil, fmt.Errorf("ERROR: hash could no be cracked")
}

func Prepare(wordlistPath string) {
    contents, err := os.ReadFile(wordlistPath)
    if err != nil {
        log.Fatal("Could not locate wordlist file")
    }
    wordlist = bytes.Split(contents, []byte("\n"))
}

func BruteHash(hash string, hashType Algorithm) BruteforceResult {
    if len(wordlist) == 0 {
        panic("Wordlist is empty")
    }

    if hashType != SHA256 {
        return BruteforceResult{hash, nil, fmt.Errorf("ERROR: unsupported hash")}
    } else if len(hash) != 64 {
        return BruteforceResult{hash, nil, fmt.Errorf("ERROR: hash cannot be bruteforced")}
    }

    hashBytes, err := hex.DecodeString(hash)
    if err != nil {
    return BruteforceResult{hash , nil, fmt.Errorf("ERROR: %v", err)}
}
    bruter := &Bruter{
        algorithm: hashType,
        hash: hashBytes,
    }
    result, err := bruter.Brute()
    return BruteforceResult{hash, result, err}
}
func IsError(result BruteforceResult) bool {
    return result.Err != nil
}
