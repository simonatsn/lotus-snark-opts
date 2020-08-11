package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/filecoin-project/sector-storage/ffiwrapper"
	"github.com/filecoin-project/specs-actors/actors/abi"
)

func BenchmarkVerifyWinning(b *testing.B) {
	buf, err := ioutil.ReadFile("tests/winning.prf")
	if err != nil {
		panic(err)
	}

	var info abi.WinningPoStVerifyInfo
	err = info.UnmarshalCBOR(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		ok, err := ffiwrapper.ProofVerifier.VerifyWinningPoSt(context.TODO(),
			info)
		if err != nil || !ok {
			b.Fatal("verify failed")
		}
	}
}

func BenchmarkVerifySeal(b *testing.B) {
	buf, err := ioutil.ReadFile("tests/seal.prf")
	if err != nil {
		panic(err)
	}

	var info abi.SealVerifyInfo
	err = info.UnmarshalCBOR(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		ok, err := ffiwrapper.ProofVerifier.VerifySeal(info)
		if err != nil || !ok {
			b.Fatal("verify failed")
		}
	}
}

func BenchmarkVerifyWindow(b *testing.B) {
	// Read in the proofs
	dir := "tests/large_window"
	proofs := make([]abi.WindowPoStVerifyInfo, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		buf, err := ioutil.ReadFile(dir + "/" + f.Name())
		if err != nil {
			panic(err)
		}

		var info abi.WindowPoStVerifyInfo
		err = info.UnmarshalCBOR(bytes.NewReader(buf))
		if err != nil {
			panic(err)
		}

		proofs = append(proofs, info)
	}

	// Verify all proofs
	for i := 0; i < b.N; i++ {
		for _, proof := range proofs {
			ok, err := ffiwrapper.ProofVerifier.VerifyWindowPoSt(context.TODO(),
				proof)
			if err != nil || !ok {
				b.Fatal("verify failed")
			}
		}
	}
}

func TestMain(m *testing.M) {
	// Warmup
	// This is particularly important for window since the first verification
	// takes a long time and we want to see the hot verification time

	buf, err := ioutil.ReadFile("tests/window.prf")
	if err != nil {
		panic(err)
	}

	var info abi.WindowPoStVerifyInfo
	err = info.UnmarshalCBOR(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	ok, err := ffiwrapper.ProofVerifier.VerifyWindowPoSt(context.TODO(), info)
	if err != nil || !ok {
		fmt.Println("verify failed")
		os.Exit(1)
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
