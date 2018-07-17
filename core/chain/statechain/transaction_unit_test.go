package statechain

import (
	"reflect"
	"testing"
)

var (
	//payload = map[string]interface{}{
	//"foo":    "bar",
	//"foobar": 5,
	//}
	payload = `{"foo":"bar"}`
	txHash  = "0xHash"
	sig     = &TxSig{
		R: " 0x0",
		S: "0x1",
	}
	props1 = &TransactionProps{
		TxHash:    nil,
		ImageHash: "0x1",
		Method:    "0x2",
		Payload:   nil,
		From:      "0x3",
		Sig:       nil,
	}
	props2 = &TransactionProps{
		TxHash:    &txHash,
		ImageHash: "0x1",
		Method:    "0x2",
		Payload:   []byte(payload),
		From:      "0x3",
		Sig:       sig,
	}
)

func TestSerializeDeserializeTransaction(t *testing.T) {
	t.Parallel()

	inputs := []*Transaction{
		NewTransaction(props1),
		NewTransaction(props2),
	}

	for idx, input := range inputs {
		bytes, err := input.Serialize()
		if err != nil {
			t.Errorf("test %d failed serialization\n%v", idx+1, err)
		}

		tx := new(Transaction)
		if err := tx.Deserialize(bytes); err != nil {
			t.Errorf("test %d failed deserialization\n%v", idx+1, err)
		}
		if tx == nil {
			t.Errorf("test %d failed\nnil tx", idx+1)
		}

		if input.props.ImageHash != tx.props.ImageHash || input.props.Method != tx.props.Method || input.props.From != tx.props.From {
			t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *input, *tx)
		}

		if !reflect.DeepEqual(input.props.Payload, tx.props.Payload) {
			t.Errorf("test %d failed\nexpected payload type %T: %v\n received payload %T: %v", idx+1, input.props.Payload, input.props.Payload, tx.props.Payload, tx.props.Payload)
		}

		if !reflect.DeepEqual(input.props.TxHash, tx.props.TxHash) {
			t.Errorf("test %d failed\nexpected txHash: %v\n received txHash: %v", idx+1, input.props.TxHash, tx.props.TxHash)
		}
	}
}

func TestSerializeDeserializeStringTransaction(t *testing.T) {
	t.Parallel()

	inputs := []*Transaction{
		NewTransaction(props1),
		NewTransaction(props2),
	}

	for idx, input := range inputs {
		str, err := input.SerializeString()
		if err != nil {
			t.Errorf("test %d failed serialization\n%v", idx+1, err)
		}

		tx := new(Transaction)
		if err := tx.DeserializeString(str); err != nil {
			t.Errorf("test %d failed deserialization\n%v", idx+1, err)
		}
		if tx == nil {
			t.Errorf("test %d failed\nnil tx", idx+1)
		}

		if input.props.ImageHash != tx.props.ImageHash || input.props.Method != tx.props.Method || input.props.From != tx.props.From {
			t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *input, *tx)
		}

		if !reflect.DeepEqual(input.props.Payload, tx.props.Payload) {
			t.Errorf("test %d failed\nexpected payload type %T: %v\n received payload %T: %v", idx+1, input.props.Payload, input.props.Payload, tx.props.Payload, tx.props.Payload)
		}

		if !reflect.DeepEqual(input.props.TxHash, tx.props.TxHash) {
			t.Errorf("test %d failed\nexpected txHash: %v\n received txHash: %v", idx+1, input.props.TxHash, tx.props.TxHash)
		}
	}
}
