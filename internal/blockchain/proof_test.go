package blockchain

import "testing"

func TestProof(t *testing.T) {
	bchain := InitBlockchain()
	blocks := []string{"A", "B", "C"}

	for _, v := range blocks {
		bchain.AddBlock(v)
	}

	if len(bchain.Blocks) != len(blocks)+1 {
		t.Fatalf("expected %d blocks, got %d", len(blocks)+1, len(bchain.Blocks))
	}

	if got := string(bchain.Blocks[0].Data); got != "Genesis" {
		t.Errorf("genesis block data mismatch: expected %q, got %q", "Genesis", got)
	}

	for i, expected := range blocks {
		got := string(bchain.Blocks[i+1].Data)
		if got != expected {
			t.Errorf("block %d data mismatch: expected %q, got %q", i+1, expected, got)
		}
	}


	for i, v := range bchain.Blocks {
		proof := NewProof(v)

		if !proof.Validate() {
			t.Errorf("invalid proof for block %d", i)
		}

	}
}