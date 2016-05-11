/*
   Copyright 2016 Continusec Pty Ltd

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package continusec

import (
	"fmt"
	"testing"
)

/*
	client := NewClient("7981306761429961588", "c9fc80d4e19ddbf01a4e6b5277a29e1bffa88fe047af9d0b9b36de536f85c2c6")
	log := client.VerifiableLog("fdsfasdfas")

create

	log := client.VerifiableLog("fdsfasdfas")

	err := log.Create()
	if err != nil {
		if err != ErrObjectConflict {
			t.Fatal(err)
		}
	}


add

	_, err = log.Add(&RawDataEntry{RawBytes: []byte("foo")})
	if err != nil {
		t.Fatal(err)
	}

	_, err = log.Add(&JsonEntry{JsonBytes: []byte("{\"name\":\"adam\",\"ssn\":123.45}")})
	if err != nil {
		t.Fatal(err)
	}

	_, err = log.Add(&RedactableJsonEntry{JsonBytes: []byte("{\"name\":\"adam\",\"ssn\":123.45}")})
	if err != nil {
		t.Fatal(err)
	}

block


	addResp, err := log.Add(&RawDataEntry{RawBytes: []byte("foo")})
	if err != nil {
		t.Fatal(err)
	}

	head, err := log.BlockUntilPresent(addResp)
	if err != nil {
		t.Fatal(err)
	}

check consistency


	h2, err := log.FetchVerifiedTreeHead(h1)
	if err != nil {
		t.Fatal(err)
	}

prove inclusion

	head, err := log.TreeHead(Head)
	if err != nil {
		t.Fatal(err)
	}

	inclProof, err := log.InclusionProof(head.TreeSize, &RawDataEntry{RawBytes: []byte("foo")})
	if err != nil {
		t.Fatal(err)
	}

	err = inclProof.Verify(head)
	if err != nil {
		t.Fatal(err)
	}

prove inclusion where supplied

	inclProof, err := log.InclusionProof(1, &RawDataEntry{RawBytes: []byte("foo")})
	if err != nil {
		t.Fatal(err)
	}

	h2, err := log.TreeHead(Head)
	if err != nil {
		t.Fatal(err)
	}

	h1, err := log.VerifySuppliedInclusionProof(h2, inclProof)
	if err != nil {
		t.Fatal(err)
	}




*/

func TestStuff(t *testing.T) {
	client := NewClient("7981306761429961588", "c9fc80d4e19ddbf01a4e6b5277a29e1bffa88fe047af9d0b9b36de536f85c2c6")
	log := client.VerifiableLog("fdsfasdfas")

	/*	inclProof, err := log.InclusionProof(1, &RawDataEntry{RawBytes: []byte("foo")})
		if err != nil {
			t.Fatal(err)
		}

		h2, err := log.TreeHead(Head)
		if err != nil {
			t.Fatal(err)
		}

		h1, err := log.VerifySuppliedInclusionProof(h2, inclProof)
		if err != nil {
			t.Fatal(err)
		}*/

	head, err := log.FetchVerifiedTreeHead(ZeroLogTreeHead)
	if err != nil {
		t.Fatal(err)
	}

	err = log.FetchAndAuditLogEntries(ZeroLogTreeHead, head, JsonEntryFactory, func(idx int64, entry VerifiableEntry) error {
		dat, err := entry.Data()
		if err != nil {
			return err
		}
		t.Log(fmt.Sprintf("idx: %d, len: %d", idx, len(dat)))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	/*

		inclProof, err := log.InclusionProof(2, &RawDataEntry{RawBytes: []byte("foo")})
		if err != nil {
			t.Fatal(err)
		}


		head, err := log.TreeHead(Head)
		if err != nil {
			t.Fatal(err)
		}

		inclusionHead, err := log.TreeHead(inclProof.TreeSize)
		if err != nil {
			t.Fatal(err)
		}

		if inclusionHead.TreeSize < head.TreeSize {



		err = inclProof.Verify(head)
		if err != nil {
			t.Fatal(err)
		}
	*/

	/*	err := log.Create()
		if err != nil {
			if err != ErrObjectConflict {
				t.Fatal(err)
			}
		}

		_, err = log.Add(&JsonEntry{JsonBytes: []byte("{\"name\":\"adam\",\"ssn\":123.45}")})
		if err != nil {
			t.Fatal(err)
		}

		addResp, err := log.Add(&RedactableJsonEntry{JsonBytes: []byte("{\"name\":\"adam\",\"ssn\":123.45}")})
		if err != nil {
			t.Fatal(err)
		}

		_, err = log.Add(&RawDataEntry{RawBytes: []byte("foo")})
		if err != nil {
			t.Fatal(err)
		}

		head, err := log.BlockUntilPresent(addResp)
		if err != nil {
			t.Fatal(err)
		}

		inclProof, err := log.InclusionProof(head.TreeSize, addResp)
		if err != nil {
			t.Fatal(err)
		}

		err = inclProof.Verify(head)
		if err != nil {
			t.Fatal(err)
		}*/

	/*t.Log("hello")

	client := NewClient("7981306761429961588", "c9fc80d4e19ddbf01a4e6b5277a29e1bffa88fe047af9d0b9b36de536f85c2c6")
	log := client.VerifiableLog("gotest")

		err := log.Create()
		if err != nil {
			t.Log("Err:", err)
		}
		_, err = log.Add([]byte("foo"))
		if err != nil {
			t.Log("Err:", err)
		}
		_, err = log.AddJson([]byte(`{"name": "ado", "ssn": 123.45}`))
		if err != nil {
			t.Log("Err:", err)
		}
		_, err = log.AddRedactibleJson([]byte(`{"name": "ado", "ssn": 123.45}`))
		if err != nil {
			t.Log("Err:", err)
		}


	treeSize, rootHash, err := log.TreeHash(Head)
	if err != nil {
		t.Log("Err:", err)
	}
	t.Log("Size", treeSize)

	data, err := log.GetEntry(0)
	if err != nil {
		t.Log("Err:", err)
	}
	t.Log(string(data.Data))

	d1, err := log.GetJsonEntry(1)
	if err != nil {
		t.Log("Err:", err)
	}
	t.Log(string(d1.Data))

	d1, err = log.GetJsonEntry(2)
	if err != nil {
		t.Log("Err:", err)
	}
	t.Log(string(d1.Data))

	d2, err := log.GetRedactedJsonEntry(2)
	if err != nil {
		t.Log("Err:", err)
	}
	x, err := d2.ShedRedacted()
	if err != nil {
		t.Log("Err:", err)
	}
	t.Log(string(x))

	bs, err := d2.BytesForHash()
	if err != nil {
		t.Log("Err:", err)
	}
	mtl := LeafMerkleTreeHash(bs)
	leafIdx, auditPath, err := log.InclusionProof(treeSize, mtl)
	if err != nil {
		t.Log("Err:", err)
	}

	t.Log("Verify result:", VerifyLogInclusionProof(leafIdx, treeSize, mtl, rootHash, auditPath))*/
}
