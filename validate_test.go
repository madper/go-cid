package cid

import (
	"testing"

	mh "github.com/multiformats/go-multihash"
)

func TestValidateCids(t *testing.T) {
	assertTrue := func(v bool) {
		t.Helper()
		if !v {
			t.Fatal("expected success")
		}
	}
	assertFalse := func(v bool) {
		t.Helper()
		if v {
			t.Fatal("expected failure")
		}
	}

	assertTrue(IsGoodHash(mh.SHA2_256))
	assertTrue(IsGoodHash(mh.BLAKE2B_MIN + 32))
	assertTrue(IsGoodHash(mh.DBL_SHA2_256))
	assertTrue(IsGoodHash(mh.KECCAK_256))
	assertTrue(IsGoodHash(mh.SHA3))

	assertTrue(IsGoodHash(mh.SHA1))

	assertFalse(IsGoodHash(mh.BLAKE2B_MIN + 5))

	mhcid := func(code uint64, length int) *Cid {
		c := &Cid{
			version: 1,
			codec:   DagCBOR,
		}
		mhash, err := mh.Sum([]byte{}, code, length)
		if err != nil {
			t.Fatal(err)
		}
		c.hash = mhash
		return c
	}

	cases := []struct {
		cid *Cid
		err error
	}{
		{mhcid(mh.SHA2_256, 32), nil},
		{mhcid(mh.SHA2_256, 16), ErrBelowMinimumHashLength},
		{mhcid(mh.MURMUR3, 4), ErrPossiblyInsecureHashFunction},
	}

	for i, cas := range cases {
		if ValidateCid(cas.cid) != cas.err {
			t.Errorf("wrong result in case of %s (index %d). Expected: %s, got %s",
				cas.cid, i, cas.err, ValidateCid(cas.cid))
		}
	}

}
