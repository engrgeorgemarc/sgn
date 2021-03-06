// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ledgerstruct

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LedgerStructABI is the input ABI used to generate the binding from.
const LedgerStructABI = "[]"

// LedgerStructBin is the compiled bytecode used for deploying new contracts.
var LedgerStructBin = "0x60556023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a72305820528a1e514b356516bd7b025fcfaafaf6a807d68d0872a6f1d8966c046ff7457164736f6c634300050a0032"

// DeployLedgerStruct deploys a new Ethereum contract, binding an instance of LedgerStruct to it.
func DeployLedgerStruct(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LedgerStruct, error) {
	parsed, err := abi.JSON(strings.NewReader(LedgerStructABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LedgerStructBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LedgerStruct{LedgerStructCaller: LedgerStructCaller{contract: contract}, LedgerStructTransactor: LedgerStructTransactor{contract: contract}, LedgerStructFilterer: LedgerStructFilterer{contract: contract}}, nil
}

// LedgerStruct is an auto generated Go binding around an Ethereum contract.
type LedgerStruct struct {
	LedgerStructCaller     // Read-only binding to the contract
	LedgerStructTransactor // Write-only binding to the contract
	LedgerStructFilterer   // Log filterer for contract events
}

// LedgerStructCaller is an auto generated read-only Go binding around an Ethereum contract.
type LedgerStructCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerStructTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LedgerStructTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerStructFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LedgerStructFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerStructSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LedgerStructSession struct {
	Contract     *LedgerStruct     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LedgerStructCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LedgerStructCallerSession struct {
	Contract *LedgerStructCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// LedgerStructTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LedgerStructTransactorSession struct {
	Contract     *LedgerStructTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// LedgerStructRaw is an auto generated low-level Go binding around an Ethereum contract.
type LedgerStructRaw struct {
	Contract *LedgerStruct // Generic contract binding to access the raw methods on
}

// LedgerStructCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LedgerStructCallerRaw struct {
	Contract *LedgerStructCaller // Generic read-only contract binding to access the raw methods on
}

// LedgerStructTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LedgerStructTransactorRaw struct {
	Contract *LedgerStructTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLedgerStruct creates a new instance of LedgerStruct, bound to a specific deployed contract.
func NewLedgerStruct(address common.Address, backend bind.ContractBackend) (*LedgerStruct, error) {
	contract, err := bindLedgerStruct(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LedgerStruct{LedgerStructCaller: LedgerStructCaller{contract: contract}, LedgerStructTransactor: LedgerStructTransactor{contract: contract}, LedgerStructFilterer: LedgerStructFilterer{contract: contract}}, nil
}

// NewLedgerStructCaller creates a new read-only instance of LedgerStruct, bound to a specific deployed contract.
func NewLedgerStructCaller(address common.Address, caller bind.ContractCaller) (*LedgerStructCaller, error) {
	contract, err := bindLedgerStruct(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LedgerStructCaller{contract: contract}, nil
}

// NewLedgerStructTransactor creates a new write-only instance of LedgerStruct, bound to a specific deployed contract.
func NewLedgerStructTransactor(address common.Address, transactor bind.ContractTransactor) (*LedgerStructTransactor, error) {
	contract, err := bindLedgerStruct(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LedgerStructTransactor{contract: contract}, nil
}

// NewLedgerStructFilterer creates a new log filterer instance of LedgerStruct, bound to a specific deployed contract.
func NewLedgerStructFilterer(address common.Address, filterer bind.ContractFilterer) (*LedgerStructFilterer, error) {
	contract, err := bindLedgerStruct(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LedgerStructFilterer{contract: contract}, nil
}

// bindLedgerStruct binds a generic wrapper to an already deployed contract.
func bindLedgerStruct(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LedgerStructABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LedgerStruct *LedgerStructRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LedgerStruct.Contract.LedgerStructCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LedgerStruct *LedgerStructRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LedgerStruct.Contract.LedgerStructTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LedgerStruct *LedgerStructRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LedgerStruct.Contract.LedgerStructTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LedgerStruct *LedgerStructCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LedgerStruct.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LedgerStruct *LedgerStructTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LedgerStruct.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LedgerStruct *LedgerStructTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LedgerStruct.Contract.contract.Transact(opts, method, params...)
}
