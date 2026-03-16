package network

import (
	. "project/orderHandler"
	. "project/shared"
)

type NetworkInterface struct {
	ConnectedNodesUpdateCh chan<- NodeIdSet
	WorldViewMergeCh       chan<- SyncView
	WorldViewReqCh         chan chan WorldView
}
