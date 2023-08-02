package main

import (
	"bytes"
	"log"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/dgraph-io/badger/v3"
)

type KVStoreApplication struct {
	db           *badger.DB
	onGoingBlock *badger.Txn
}

var _ abcitypes.Application = (*KVStoreApplication)(nil)

func NewKVStoreApplication(db *badger.DB) *KVStoreApplication {
	return &KVStoreApplication{
		db: db,
	}
}

var (
	OK      = abcitypes.CodeTypeOK // 0 indicates everything ok
	Invalid = uint32(1)            // any non-zero code indicates an error
)

func (app *KVStoreApplication) isValid(tx []byte) bool {
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return false
	}
	return true
}

func (app *KVStoreApplication) Info(info abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (app *KVStoreApplication) Query(query abcitypes.RequestQuery) abcitypes.ResponseQuery {
	resp := abcitypes.ResponseQuery{Key: query.Data}

	dbErr := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(query.Data)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
			resp.Log = "key does not exist"
			return nil
		}

		return item.Value(func(val []byte) error {
			resp.Log = "exists"
			resp.Value = val
			return nil
		})
	})
	if dbErr != nil {
		log.Panicf("Error reading database, unable to execute query: %v", dbErr)
	}
	return resp
}

func (app *KVStoreApplication) CheckTx(tx abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	if !app.isValid(tx.Tx) {
		return abcitypes.ResponseCheckTx{Code: Invalid}
	}
	return abcitypes.ResponseCheckTx{Code: OK}
}

func (app *KVStoreApplication) InitChain(chain abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *KVStoreApplication) PrepareProposal(proposal abcitypes.RequestPrepareProposal) abcitypes.ResponsePrepareProposal {
	return abcitypes.ResponsePrepareProposal{Txs: proposal.Txs}
}

func (app *KVStoreApplication) ProcessProposal(proposal abcitypes.RequestProcessProposal) abcitypes.ResponseProcessProposal {
	return abcitypes.ResponseProcessProposal{Status: abcitypes.ResponseProcessProposal_ACCEPT}
}

func (app *KVStoreApplication) BeginBlock(_ abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	// what if we don't commit the previous block?
	app.onGoingBlock = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (app *KVStoreApplication) DeliverTx(request abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	if !app.isValid(request.Tx) {
		return abcitypes.ResponseDeliverTx{Code: Invalid}
	}

	parts := bytes.SplitN(request.Tx, []byte("="), 2)
	key, value := parts[0], parts[1]

	if err := app.onGoingBlock.Set(key, value); err != nil {
		log.Panicf("Error writing to database, unable to execute request: %v", err)
	}

	return abcitypes.ResponseDeliverTx{Code: OK}
}

func (app *KVStoreApplication) EndBlock(block abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (app *KVStoreApplication) Commit() abcitypes.ResponseCommit {
	if err := app.onGoingBlock.Commit(); err != nil {
		log.Panicf("Error writing to database, unable to commit block: %v", err)
	}
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *KVStoreApplication) ListSnapshots(snapshots abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

func (app *KVStoreApplication) OfferSnapshot(snapshot abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

func (app *KVStoreApplication) LoadSnapshotChunk(chunk abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

func (app *KVStoreApplication) ApplySnapshotChunk(chunk abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}
