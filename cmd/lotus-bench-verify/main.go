package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/filecoin-project/sector-storage/ffiwrapper"
	"github.com/filecoin-project/specs-actors/actors/abi"
)

var repetitions = 1

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func verifySeal(fname string) bool {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	var info abi.SealVerifyInfo
	err = info.UnmarshalCBOR(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	ok, err := ffiwrapper.ProofVerifier.VerifySeal(info)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("Seal BAD!")
		return false
	}
	return true
}

func verifyWindow(fname string) bool {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	var info abi.WindowPoStVerifyInfo
	err = info.UnmarshalCBOR(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	ok, err := ffiwrapper.ProofVerifier.VerifyWindowPoSt(context.TODO(), info)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("Window BAD!")
		return false
	}
	return true
}

func verifyWinning(fname string) bool {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	var info abi.WinningPoStVerifyInfo
	err = info.UnmarshalCBOR(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	ok, err := ffiwrapper.ProofVerifier.VerifyWinningPoSt(context.TODO(), info)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("Winning BAD!")
		return false
	}
	return true
}

// Verify all the proofs in the given directory that match the proofType
func verifyAll(proofType string, verify func(string) bool) bool {
	files, err := ioutil.ReadDir("tests")
	if err != nil {
		log.Fatal(err)
	}

	// Don't time the first proof
	warmedUp := false

	for _, f := range files {
		if strings.HasPrefix(f.Name(), proofType) {

			if !warmedUp {
				warmedUp = true
				if !verify("tests/" + f.Name()) {
					fmt.Println("Verify failed for", f.Name())
					return false
				}
			}

			for i := 0; i < repetitions; i++ {
				start := time.Now()
				if !verify("tests/" + f.Name()) {
					fmt.Println("Verify failed for", f.Name())
					return false
				}
				timeTrack(start, "Verify "+f.Name())
			}
		}
	}
	return true
}

func main() {
	ok := true
	if !verifyWindow("tests/window.prf") {
		ok = false
	}
	start := time.Now()
	if !verifyWindow("tests/window.prf") {
		ok = false
	}
	timeTrack(start, "Verify window")

	if !verifySeal("tests/seal.prf") {
		ok = false
	}
	start = time.Now()
	if !verifySeal("tests/seal.prf") {
		ok = false
	}
	timeTrack(start, "Verify seal")

	if !verifyWinning("tests/winning.prf") {
		ok = false
	}
	start = time.Now()
	if !verifyWinning("tests/winning.prf") {
		ok = false
	}
	timeTrack(start, "Verify winning")

	// if !verifyAll("seal", verifySeal) {
	// 	ok = false
	// }
	// if !verifyAll("window", verifyWindow) {
	// 	ok = false
	// }
	// if !verifyAll("winning", verifyWinning) {
	// 	ok = false
	// }
	if !ok {
		fmt.Println("Errors encountered")
	} else {
		fmt.Println("Success!")
	}
}
