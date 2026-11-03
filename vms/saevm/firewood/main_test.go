// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package firewood

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const fixedRapidSeed = "1"

// pinRapidSeed pins the seed only when neither -rapid.seed nor shuffling asks for variety.
func pinRapidSeed(seedSet bool, shuffle string) (seed string, pin bool) {
	if seedSet || shuffle != "off" {
		return "", false
	}
	return fixedRapidSeed, true
}

func TestMain(m *testing.M) {
	flag.Parse()
	seedSet := false
	flag.Visit(func(f *flag.Flag) { seedSet = seedSet || f.Name == "rapid.seed" })
	if seed, pin := pinRapidSeed(seedSet, flag.Lookup("test.shuffle").Value.String()); pin {
		if err := flag.Set("rapid.seed", seed); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

func TestPinRapidSeed(t *testing.T) {
	tests := []struct {
		name    string
		seedSet bool
		shuffle string
		wantPin bool
	}{
		{"explicit seed", true, "off", false},
		{"shuffle on", false, "on", false},
		{"shuffle seed", false, "12345", false},
		{"defaults", false, "off", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed, pin := pinRapidSeed(tt.seedSet, tt.shuffle)
			require.Equal(t, tt.wantPin, pin, "pinRapidSeed(%v, %q) pin", tt.seedSet, tt.shuffle)
			if tt.wantPin {
				require.Equal(t, fixedRapidSeed, seed, "pinRapidSeed() seed")
			}
		})
	}
}
