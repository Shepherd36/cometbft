package types

import (
	"os"
	"testing"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/internal/test"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var config *cfg.Config // NOTE: must be reset for each _test.go file

func TestMain(m *testing.M) {
	config = test.ResetTestRoot("consensus_height_vote_set_test")
	code := m.Run()
	os.RemoveAll(config.RootDir)
	os.Exit(code)
}

func TestPeerCatchupRounds(t *testing.T) {
	valSet, privVals := types.RandValidatorSet(10, 1)

	hvs := NewExtendedHeightVoteSet(test.DefaultTestChainID, 1, valSet)

	vote999_0 := makeVoteHR(t, 1, 0, 999, privVals)
	added, err := hvs.AddVote(vote999_0, "peer1")
	if !added || err != nil {
		t.Error("Expected to successfully add vote from peer", added, err)
	}

	vote1000_0 := makeVoteHR(t, 1, 0, 1000, privVals)
	added, err = hvs.AddVote(vote1000_0, "peer1")
	if !added || err != nil {
		t.Error("Expected to successfully add vote from peer", added, err)
	}

	vote1001_0 := makeVoteHR(t, 1, 0, 1001, privVals)
	added, err = hvs.AddVote(vote1001_0, "peer1")
	if err != ErrGotVoteFromUnwantedRound {
		t.Errorf("expected GotVoteFromUnwantedRoundError, but got %v", err)
	}
	if added {
		t.Error("Expected to *not* add vote from peer, too many catchup rounds.")
	}

	added, err = hvs.AddVote(vote1001_0, "peer2")
	if !added || err != nil {
		t.Error("Expected to successfully add vote from another peer")
	}

}

func makeVoteHR(
	t *testing.T,
	height int64,
	valIndex,
	round int32,
	privVals []types.PrivValidator,
) *types.Vote {
	privVal := privVals[valIndex]
	randBytes := tmrand.Bytes(tmhash.Size)

	vote, err := types.MakeVote(
		privVal,
		test.DefaultTestChainID,
		valIndex,
		height,
		round,
		tmproto.PrecommitType,
		types.BlockID{Hash: randBytes, PartSetHeader: types.PartSetHeader{}},
		tmtime.Now(),
	)
	if err != nil {
		panic(err)
	}

	return vote
}
