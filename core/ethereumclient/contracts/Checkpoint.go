// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package checkpoint

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// CheckpointABI is the input ABI used to generate the binding from.
const CheckpointABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkpoints\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkpointId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"checkpointId\",\"type\":\"uint256\"}],\"name\":\"LogCheckpoint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"root\",\"type\":\"uint256\"}],\"name\":\"checkpoint\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"bytes32\"}],\"name\":\"verify\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// Checkpoint is an auto generated Go binding around an Ethereum contract.
type Checkpoint struct {
	CheckpointCaller     // Read-only binding to the contract
	CheckpointTransactor // Write-only binding to the contract
	CheckpointFilterer   // Log filterer for contract events
}

// CheckpointCaller is an auto generated read-only Go binding around an Ethereum contract.
type CheckpointCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckpointTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CheckpointTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckpointFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CheckpointFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckpointSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CheckpointSession struct {
	Contract     *Checkpoint       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CheckpointCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CheckpointCallerSession struct {
	Contract *CheckpointCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CheckpointTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CheckpointTransactorSession struct {
	Contract     *CheckpointTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CheckpointRaw is an auto generated low-level Go binding around an Ethereum contract.
type CheckpointRaw struct {
	Contract *Checkpoint // Generic contract binding to access the raw methods on
}

// CheckpointCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CheckpointCallerRaw struct {
	Contract *CheckpointCaller // Generic read-only contract binding to access the raw methods on
}

// CheckpointTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CheckpointTransactorRaw struct {
	Contract *CheckpointTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCheckpoint creates a new instance of Checkpoint, bound to a specific deployed contract.
func NewCheckpoint(address common.Address, backend bind.ContractBackend) (*Checkpoint, error) {
	contract, err := bindCheckpoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Checkpoint{CheckpointCaller: CheckpointCaller{contract: contract}, CheckpointTransactor: CheckpointTransactor{contract: contract}, CheckpointFilterer: CheckpointFilterer{contract: contract}}, nil
}

// NewCheckpointCaller creates a new read-only instance of Checkpoint, bound to a specific deployed contract.
func NewCheckpointCaller(address common.Address, caller bind.ContractCaller) (*CheckpointCaller, error) {
	contract, err := bindCheckpoint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CheckpointCaller{contract: contract}, nil
}

// NewCheckpointTransactor creates a new write-only instance of Checkpoint, bound to a specific deployed contract.
func NewCheckpointTransactor(address common.Address, transactor bind.ContractTransactor) (*CheckpointTransactor, error) {
	contract, err := bindCheckpoint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CheckpointTransactor{contract: contract}, nil
}

// NewCheckpointFilterer creates a new log filterer instance of Checkpoint, bound to a specific deployed contract.
func NewCheckpointFilterer(address common.Address, filterer bind.ContractFilterer) (*CheckpointFilterer, error) {
	contract, err := bindCheckpoint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CheckpointFilterer{contract: contract}, nil
}

// bindCheckpoint binds a generic wrapper to an already deployed contract.
func bindCheckpoint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CheckpointABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Checkpoint *CheckpointRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Checkpoint.Contract.CheckpointCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Checkpoint *CheckpointRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Checkpoint.Contract.CheckpointTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Checkpoint *CheckpointRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Checkpoint.Contract.CheckpointTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Checkpoint *CheckpointCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Checkpoint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Checkpoint *CheckpointTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Checkpoint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Checkpoint *CheckpointTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Checkpoint.Contract.contract.Transact(opts, method, params...)
}

// CheckpointId is a free data retrieval call binding the contract method 0xc18abbeb.
//
// Solidity: function checkpointId() constant returns(uint256)
func (_Checkpoint *CheckpointCaller) CheckpointId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Checkpoint.contract.Call(opts, out, "checkpointId")
	return *ret0, err
}

// CheckpointId is a free data retrieval call binding the contract method 0xc18abbeb.
//
// Solidity: function checkpointId() constant returns(uint256)
func (_Checkpoint *CheckpointSession) CheckpointId() (*big.Int, error) {
	return _Checkpoint.Contract.CheckpointId(&_Checkpoint.CallOpts)
}

// CheckpointId is a free data retrieval call binding the contract method 0xc18abbeb.
//
// Solidity: function checkpointId() constant returns(uint256)
func (_Checkpoint *CheckpointCallerSession) CheckpointId() (*big.Int, error) {
	return _Checkpoint.Contract.CheckpointId(&_Checkpoint.CallOpts)
}

// Checkpoints is a free data retrieval call binding the contract method 0xb8a24252.
//
// Solidity: function checkpoints( uint256) constant returns(uint256)
func (_Checkpoint *CheckpointCaller) Checkpoints(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Checkpoint.contract.Call(opts, out, "checkpoints", arg0)
	return *ret0, err
}

// Checkpoints is a free data retrieval call binding the contract method 0xb8a24252.
//
// Solidity: function checkpoints( uint256) constant returns(uint256)
func (_Checkpoint *CheckpointSession) Checkpoints(arg0 *big.Int) (*big.Int, error) {
	return _Checkpoint.Contract.Checkpoints(&_Checkpoint.CallOpts, arg0)
}

// Checkpoints is a free data retrieval call binding the contract method 0xb8a24252.
//
// Solidity: function checkpoints( uint256) constant returns(uint256)
func (_Checkpoint *CheckpointCallerSession) Checkpoints(arg0 *big.Int) (*big.Int, error) {
	return _Checkpoint.Contract.Checkpoints(&_Checkpoint.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Checkpoint *CheckpointCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Checkpoint.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Checkpoint *CheckpointSession) IsOwner() (bool, error) {
	return _Checkpoint.Contract.IsOwner(&_Checkpoint.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Checkpoint *CheckpointCallerSession) IsOwner() (bool, error) {
	return _Checkpoint.Contract.IsOwner(&_Checkpoint.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Checkpoint *CheckpointCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Checkpoint.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Checkpoint *CheckpointSession) Owner() (common.Address, error) {
	return _Checkpoint.Contract.Owner(&_Checkpoint.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Checkpoint *CheckpointCallerSession) Owner() (common.Address, error) {
	return _Checkpoint.Contract.Owner(&_Checkpoint.CallOpts)
}

// Verify is a free data retrieval call binding the contract method 0x5a9a49c7.
//
// Solidity: function verify(proof bytes32[], root bytes32, leaf bytes32) constant returns(bool)
func (_Checkpoint *CheckpointCaller) Verify(opts *bind.CallOpts, proof [][32]byte, root [32]byte, leaf [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Checkpoint.contract.Call(opts, out, "verify", proof, root, leaf)
	return *ret0, err
}

// Verify is a free data retrieval call binding the contract method 0x5a9a49c7.
//
// Solidity: function verify(proof bytes32[], root bytes32, leaf bytes32) constant returns(bool)
func (_Checkpoint *CheckpointSession) Verify(proof [][32]byte, root [32]byte, leaf [32]byte) (bool, error) {
	return _Checkpoint.Contract.Verify(&_Checkpoint.CallOpts, proof, root, leaf)
}

// Verify is a free data retrieval call binding the contract method 0x5a9a49c7.
//
// Solidity: function verify(proof bytes32[], root bytes32, leaf bytes32) constant returns(bool)
func (_Checkpoint *CheckpointCallerSession) Verify(proof [][32]byte, root [32]byte, leaf [32]byte) (bool, error) {
	return _Checkpoint.Contract.Verify(&_Checkpoint.CallOpts, proof, root, leaf)
}

// Checkpoint is a paid mutator transaction binding the contract method 0xed64bab2.
//
// Solidity: function checkpoint(root uint256) returns()
func (_Checkpoint *CheckpointTransactor) Checkpoint(opts *bind.TransactOpts, root *big.Int) (*types.Transaction, error) {
	return _Checkpoint.contract.Transact(opts, "checkpoint", root)
}

// Checkpoint is a paid mutator transaction binding the contract method 0xed64bab2.
//
// Solidity: function checkpoint(root uint256) returns()
func (_Checkpoint *CheckpointSession) Checkpoint(root *big.Int) (*types.Transaction, error) {
	return _Checkpoint.Contract.Checkpoint(&_Checkpoint.TransactOpts, root)
}

// Checkpoint is a paid mutator transaction binding the contract method 0xed64bab2.
//
// Solidity: function checkpoint(root uint256) returns()
func (_Checkpoint *CheckpointTransactorSession) Checkpoint(root *big.Int) (*types.Transaction, error) {
	return _Checkpoint.Contract.Checkpoint(&_Checkpoint.TransactOpts, root)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Checkpoint *CheckpointTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Checkpoint.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Checkpoint *CheckpointSession) RenounceOwnership() (*types.Transaction, error) {
	return _Checkpoint.Contract.RenounceOwnership(&_Checkpoint.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Checkpoint *CheckpointTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Checkpoint.Contract.RenounceOwnership(&_Checkpoint.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Checkpoint *CheckpointTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Checkpoint.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Checkpoint *CheckpointSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Checkpoint.Contract.TransferOwnership(&_Checkpoint.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Checkpoint *CheckpointTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Checkpoint.Contract.TransferOwnership(&_Checkpoint.TransactOpts, newOwner)
}

// CheckpointLogCheckpointIterator is returned from FilterLogCheckpoint and is used to iterate over the raw logs and unpacked data for LogCheckpoint events raised by the Checkpoint contract.
type CheckpointLogCheckpointIterator struct {
	Event *CheckpointLogCheckpoint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CheckpointLogCheckpointIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CheckpointLogCheckpoint)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CheckpointLogCheckpoint)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CheckpointLogCheckpointIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CheckpointLogCheckpointIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CheckpointLogCheckpoint represents a LogCheckpoint event raised by the Checkpoint contract.
type CheckpointLogCheckpoint struct {
	CheckpointId *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLogCheckpoint is a free log retrieval operation binding the contract event 0xf33e784ea02a76519c911149a08a2f29a8e2f3ee4bb824b7080c09bf635d692f.
//
// Solidity: e LogCheckpoint(checkpointId indexed uint256)
func (_Checkpoint *CheckpointFilterer) FilterLogCheckpoint(opts *bind.FilterOpts, checkpointId []*big.Int) (*CheckpointLogCheckpointIterator, error) {

	var checkpointIdRule []interface{}
	for _, checkpointIdItem := range checkpointId {
		checkpointIdRule = append(checkpointIdRule, checkpointIdItem)
	}

	logs, sub, err := _Checkpoint.contract.FilterLogs(opts, "LogCheckpoint", checkpointIdRule)
	if err != nil {
		return nil, err
	}
	return &CheckpointLogCheckpointIterator{contract: _Checkpoint.contract, event: "LogCheckpoint", logs: logs, sub: sub}, nil
}

// WatchLogCheckpoint is a free log subscription operation binding the contract event 0xf33e784ea02a76519c911149a08a2f29a8e2f3ee4bb824b7080c09bf635d692f.
//
// Solidity: e LogCheckpoint(checkpointId indexed uint256)
func (_Checkpoint *CheckpointFilterer) WatchLogCheckpoint(opts *bind.WatchOpts, sink chan<- *CheckpointLogCheckpoint, checkpointId []*big.Int) (event.Subscription, error) {

	var checkpointIdRule []interface{}
	for _, checkpointIdItem := range checkpointId {
		checkpointIdRule = append(checkpointIdRule, checkpointIdItem)
	}

	logs, sub, err := _Checkpoint.contract.WatchLogs(opts, "LogCheckpoint", checkpointIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CheckpointLogCheckpoint)
				if err := _Checkpoint.contract.UnpackLog(event, "LogCheckpoint", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// CheckpointOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Checkpoint contract.
type CheckpointOwnershipTransferredIterator struct {
	Event *CheckpointOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CheckpointOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CheckpointOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CheckpointOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CheckpointOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CheckpointOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CheckpointOwnershipTransferred represents a OwnershipTransferred event raised by the Checkpoint contract.
type CheckpointOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Checkpoint *CheckpointFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CheckpointOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Checkpoint.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CheckpointOwnershipTransferredIterator{contract: _Checkpoint.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Checkpoint *CheckpointFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CheckpointOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Checkpoint.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CheckpointOwnershipTransferred)
				if err := _Checkpoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
