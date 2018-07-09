package mainchain

import (
	"crypto/ecdsa"
	"errors"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
)

// VerifyBlock verifies a block
// TODO: check block time
// TODO: fetch and check previous block hash
func VerifyBlock(block *Block) (bool, error) {
	if block == nil {
		return false, errors.New("block is nil")
	}

	if block.props.BlockHash == nil {
		return false, errors.New("block hash is nil")
	}

	if ImageHash != block.props.ImageHash {
		return false, nil
	}

	if block.props.MinerSig == nil {
		return false, nil
	}

	difficulty, err := hexutil.DecodeUint64(block.props.Difficulty)
	if err != nil {
		return false, err
	}

	// TODO: check the difficulty is correct

	hashStr := hexutil.StripLeader(block.props.BlockHash)
	if hashStr <= difficulty {
		return false, nil
	}

	for i := 0; i < difficulty; i++ {
		if hashStr[i:i+1] != "0" {
			return false, nil
		}
	}

	// hash must verify
	tmpProps := BlockProps{
		BlockNumber:          block.props.BlockNumber,
		BlockTime:            block.props.BlockTime,
		ImageHash:            block.props.ImageHash,
		StateBlockMerkleHash: block.props.StateBlockMerkleHash,
		StateBlockHashes:     block.props.StateBlockHashes,
		PrevBlockHash:        block.props.PrevBlockHash,
		Nonce:                block.props.Nonce,
		Difficulty:           block.props.Difficulty,
		MinerAddress:         block.props.MinerAddress,
	}
	tmpBlock := Block{
		props: tmpProps,
	}
	tmpHash, err := tmpBlock.CalcHash()
	if err != nil {
		return false, err
	}
	// note: already checked for nil hash
	if *block.props.BlockHash != tmpHash {
		return false, nil
	}

	// the sig must verify
	pub, err := PubFromBlock(block)
	if err != nil {
		return false, err
	}

	// note: checked for nil sig, above
	r, err := hexutil.DecodeBigInt(block.props.MinerSig.R)
	if err != nil {
		return false, err
	}
	s, err := hexutil.DecodeBigInt(block.props.MinerSig.S)
	if err != nil {
		return false, err
	}

	c3crypto.Verify(pub, []byte(*tx.props.TxHash), r, s)

	// TODO: do in go funcs
	for _, stateblockHash := range block.props.StateBlockHashes {
		//stateblockCid, err := p2p.GetCIDByHash(stateblockHash)
		//if err != nil {
		//return false, err
		//}

		//// TODO: need to move p2p inside here
		//stateblock, err := s.props.P2P.GetStatechainBlock(stateblockCid)
		//if err != nil {
		//return false, err
		//}

		//ok, err := statechain.VerifyBlock(stateblock)
		//if err != nil {
		//return false, err
		//}
		//if !ok {
		//return false, nil
		//}
	}

	return true, nil
}

// NewFromStateBlocks ...
// TODO: everything...
func NewFromStateBlocks(stateBlocks []*statechain.Block) (*Block, error) {
	return nil, nil
}

// PubFromBlock ...
func PubFromBlock(block *mainchain.Block) (*ecdsa.PublicKey, error) {
	if block == nil {
		return nil, errors.New("block is nil")
	}

	pubStr, err := hexutil.DecodeString(tx.props.MinerAddress)
	if err != nil {
		return nil, err
	}
	pub, err := c3crypto.DeserializePublicKey([]byte(pubStr))
	if err != nil {
		return nil, err
	}
	if pub == nil {
		return nil, errors.New("invalid miner address")
	}

	return pub, nil
}
