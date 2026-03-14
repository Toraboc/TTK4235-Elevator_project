package network

import (
	. "project/shared"
	. "project/orderHandler"
)

type NetworkInterface struct {
	ConnectedNodesUpdateCh chan<- NodeIdSet
	WorldViewMergeCh       chan<- SyncView
	WorldViewReqCh         chan chan WorldView
}
