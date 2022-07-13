package whois

const (
	AddPeriod = 1 << iota
	AutoRenewPeriod
	Inactive
	Ok
	PendingCreate
	PendingDelete
	PendingRenew
	PendingRestore
	PendingTransfer
	PendingUpdate
	RedemptionPeriod
	RenewPeriod
	ServerDeleteProhibited
	ServerHold
	ServerRenewProhibited
	ServerTransferProhibited
	ServerUpdateProhibited
	TransferPeriod
	ClientDeleteProhibited
	ClientHold
	ClientRenewProhibited
	ClientTransferProhibited
	ClientUpdateProhibited
)

var (
	StatusList     = []uint{AddPeriod, AutoRenewPeriod, Inactive, Ok, PendingCreate, PendingDelete, PendingRenew, PendingRestore, PendingTransfer, PendingUpdate, RedemptionPeriod, RenewPeriod, ServerDeleteProhibited, ServerHold, ServerRenewProhibited, ServerTransferProhibited, ServerUpdateProhibited, TransferPeriod, ClientDeleteProhibited, ClientHold, ClientRenewProhibited, ClientTransferProhibited, ClientUpdateProhibited}
	StatusToString = map[uint]string{
		AddPeriod:                "addPeriod",
		AutoRenewPeriod:          "autoRenewPeriod",
		Inactive:                 "inactive",
		Ok:                       "ok",
		PendingCreate:            "pendingCreate",
		PendingDelete:            "pendingDelete",
		PendingRenew:             "pendingRenew",
		PendingRestore:           "pendingRestore",
		PendingTransfer:          "pendingTransfer",
		PendingUpdate:            "pendingUpdate",
		RedemptionPeriod:         "redemptionPeriod",
		RenewPeriod:              "renewPeriod",
		ServerDeleteProhibited:   "serverDeleteProhibited",
		ServerHold:               "serverHold",
		ServerRenewProhibited:    "serverRenewProhibited",
		ServerTransferProhibited: "serverTransferProhibited",
		ServerUpdateProhibited:   "serverUpdateProhibited",
		TransferPeriod:           "transferPeriod",
		ClientDeleteProhibited:   "clientDeleteProhibited",
		ClientHold:               "clientHold",
		ClientRenewProhibited:    "clientRenewProhibited",
		ClientTransferProhibited: "clientTransferProhibited",
		ClientUpdateProhibited:   "clientUpdateProhibited",
	}
	StatusFromString = map[string]uint{
		"addperiod":                AddPeriod,
		"autorenewperiod":          AutoRenewPeriod,
		"inactive":                 Inactive,
		"ok":                       Ok,
		"pendingcreate":            PendingCreate,
		"pendingdelete":            PendingDelete,
		"pendingrenew":             PendingRenew,
		"pendingrestore":           PendingRestore,
		"pendingtransfer":          PendingTransfer,
		"pendingupdate":            PendingUpdate,
		"redemptionperiod":         RedemptionPeriod,
		"renewperiod":              RenewPeriod,
		"serverdeleteprohibited":   ServerDeleteProhibited,
		"serverhold":               ServerHold,
		"serverrenewprohibited":    ServerRenewProhibited,
		"servertransferprohibited": ServerTransferProhibited,
		"serverupdateprohibited":   ServerUpdateProhibited,
		"transferperiod":           TransferPeriod,
		"clientdeleteprohibited":   ClientDeleteProhibited,
		"clienthold":               ClientHold,
		"clientrenewprohibited":    ClientRenewProhibited,
		"clienttransferprohibited": ClientTransferProhibited,
		"clientupdateprohibited":   ClientUpdateProhibited,
	}
)
