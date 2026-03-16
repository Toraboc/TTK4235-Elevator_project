package network

import (
	. "project/orderHandler"
	. "project/shared"
)

type NetworkInterface struct {
	ConnectedNodesUpdateCh chan<- NodeIdSet
	SyncMergeCh            chan<- SyncData
	RequestSyncCh          chan chan SyncData
}
