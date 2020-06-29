package types

import (
	"github.com/celer-network/sgn/mainchain"
)

const (
	QuerySubscription = "subscription"
	QueryRequest      = "request"
	QueryEpoch        = "epoch"
	QueryParameters   = "parameters"
)

type QuerySubscriptionParams struct {
	EthAddress string
}

func NewQuerySubscriptionParams(addr string) QuerySubscriptionParams {
	return QuerySubscriptionParams{
		EthAddress: mainchain.FormatAddrHex(addr),
	}
}

type QueryRequestParams struct {
	ChannelId []byte
	Receiver  string
}

func NewQueryRequestParams(channelId []byte, receiver string) QueryRequestParams {
	return QueryRequestParams{
		ChannelId: channelId,
		Receiver:  mainchain.FormatAddrHex(receiver),
	}
}

type QueryEpochParams struct {
	EpochId int64
}

func NewQueryEpochParams(epochId int64) QueryEpochParams {
	return QueryEpochParams{
		EpochId: epochId,
	}
}
