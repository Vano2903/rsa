package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"bufio"
	"os"
	
	"github.com/fxtlabs/primes"
)

type entity struct {
	N          int
	PrivateKey int
	PublicKey  int
}

func (ent *entity) generateKeys() error {
	found := false
	var n, e, d int

	for {
		p, err := getPrimeNumber()
		if err != nil {
			return err
		}

		q, err := getPrimeNumber()
		if err != nil {
			return err
		}

		if q == p {
			continue
		}

		fmt.Println("p:", p)
		fmt.Println("q:", q)

		n = p * q
		φ := (p - 1) * (q - 1)
		for i := φ - 1; ; i-- {
			if primes.Coprime(φ, i) {
				fmt.Println("e:", i)
				e = i
				break
			}
		}

		for i := 0; i < φ*φ; i++ {
			res := (i * e) % φ
			if res == 1 && i != e { //
				d = i
				fmt.Println("d:", d)
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	ent.N = n
	ent.PrivateKey = e
	ent.PublicKey = d

	fmt.Println("-------------------------------")
	fmt.Println("n:", ent.N)
	fmt.Println("private key:", ent.PrivateKey)
	fmt.Println("public key:", ent.PublicKey)

	return nil
}

func (ent entity) encrypt(toEncrypt []byte) []byte {
	var result []byte
	for _, b := range toEncrypt {
		result = append(result, byte((int(b)^ent.PrivateKey)%ent.N))
	}
	return result
}

func (ent entity) decrypt(toDecrypt []byte) []byte {
	var result []byte
	for _, b := range toDecrypt {
		result = append(result, byte((int(b)^ent.PrivateKey)%ent.N))
	}
	return result
}

func getPrimeNumber() (int, error) {
	n, err := rand.Prime(rand.Reader, 12)
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func main() {
	var e entity
	fmt.Println("generating keys")
	err := e.generateKeys()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-------------------------------")
	fmt.Print("message to ecrypt:")
	in := bufio.NewReader(os.Stdin)
	message, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	result := e.encrypt([]byte(message))
	fmt.Println("encrypted message:", string(result))

	result = e.decrypt(result)
	fmt.Println("decrypted message:", string(result))
}
