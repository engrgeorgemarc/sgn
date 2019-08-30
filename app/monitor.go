package app

import (
	"context"
	"log"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthMonitor struct {
	ethClient  *mainchain.EthClient
	transactor *utils.Transactor
	cdc        *codec.Codec
	started    bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, transactor *utils.Transactor, cdc *codec.Codec) *EthMonitor {
	return &EthMonitor{
		ethClient:  ethClient,
		transactor: transactor,
		cdc:        cdc,
	}
}

func (m *EthMonitor) Start() {
	if m.started {
		return
	}

	go m.monitorBlockHead()
	go m.monitorIntendSettle()
	m.started = true
}

func (m *EthMonitor) monitorBlockHead() {
	headerChan := make(chan *types.Header)
	sub, err := m.ethClient.Client.SubscribeNewHead(context.Background(), headerChan)
	if err != nil {
		log.Printf("SubscribeNewHead err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("SubscribeNewHead err", err)
		case header := <-headerChan:
			msg := global.NewMsgSyncBlock(header.Number.Uint64(), m.transactor.Key.GetAddress())
			_, err := m.transactor.BroadcastTx(msg)
			if err != nil {
				log.Printf("SyncBlock err", err)
				return
			}
		}
	}
}

func (m *EthMonitor) monitorIntendSettle() {
	intendSettleChan := make(chan *mainchain.CelerLedgerIntendSettle)
	sub, err := m.ethClient.Ledger.WatchIntendSettle(nil, intendSettleChan, [][32]byte{})
	if err != nil {
		log.Printf("WatchIntendSettle err", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WatchIntendSettle err", err)
		case intendSettle := <-intendSettleChan:
			m.handleIntendSettle(intendSettle)
		}
	}
}

func (m *EthMonitor) handleIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	request, err := subscribe.CLIQueryRequest(m.cdc, m.transactor.CliCtx, subscribe.StoreKey, intendSettle.ChannelId[:])
	if err != nil {
		log.Printf("Query request err", err)
		return
	}

	if intendSettle.SeqNums[request.PeerFromIndex].Uint64() >= request.SeqNum {
		log.Printf("Ignore the intendSettle event due to larger seqNum")
		return
	}

	tx, err := m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, request.SignedSimplexStateBytes)
	if err != nil {
		log.Printf("Tx err", err)
		return
	}
	log.Printf("Tx detail", tx)
}
