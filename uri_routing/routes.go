package uri_routing

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"GET",
		"/upload",
		Upload,
	},
	Route{
		"Upload",
		"GET",
		"/uploadpids",
		UploadPids,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	Route{
		"Canonical",
		"GET",
		"/canonical",
		Canonical,
	},
	Route{
		"ShowBlockMpt",
		"GET",
		"/showBlockMpt/{height}",
		ShowBlockMpt,
	},
	Route{
		"ShowBlock",
		"GET",
		"/showBlock/{height}",
		ShowBlock,
	},
	////////currency
	//Route{
	//	"Transaction",
	//	"POST",
	//	"/transaction", //to put in tx pool
	//	Transaction,
	//},
	Route{
		"ShowWallet",
		"GET",
		"/showWallet", //to put in tx pool
		ShowWallet,
	},
	Route{
		"ShowBalanceBook",
		"GET",
		"/showBalanceBook", //to put in tx pool
		ShowBalanceBook,
	},
	Route{
		"ShowTransactionPool",
		"GET",
		"/showTransactionPool", //to put in tx pool
		ShowTransactionPool,
	},
	Route{
		"NewClient", // used by client
		"GET",
		"/",
		NewClient,
	},
	Route{
		"ServeClient", //used by bcHolder
		"GET",
		"/client",
		ServeClient,
	},
	Route{
		"BcHolders", //called by client - return value is defined and described in BcHolders
		"GET",
		"/bcholders",
		BcHolders,
	},
	Route{
		"ShowBcHolders", //called by client
		"GET",
		"/showbcholders",
		ShowBcHolders,
	},
	Route{
		"SignUp", // used by client
		"POST",
		"/signup",
		SignUp,
	},
	Route{
		"ClientSignUp", // used by bcHolder
		"POST",
		"/clientsignup",
		ClientSignUp,
	},
	Route{
		"Login", // used by client
		"POST",
		"/login",
		Login,
	},
	Route{
		"ClientLogin", // used by bcHolder
		"POST",
		"/clientlogin",
		ClientLogin,
	},
	Route{
		"CIDSet", // used by client to set CID after correct login // NEW !!!!
		"POST",
		"/cidset",
		CIDSet,
	},
	Route{
		"CIDPage", // used by client to set CID //todo - remove - temporary solution
		"get",
		"/cidpage",
		CIDPage,
	},
	Route{
		"SetCID", // used by client to set CID //todo - remove - temporary solution
		"POST",
		"/setcid",
		SetCID,
	},
	Route{
		"TransactionForm", //used by client - submit of transaction form - handled by client
		"POST",
		"/transactionform",
		TransactionForm,
	},
	Route{
		"TransactionBeatRecv", //api of bcHolder
		"POST",
		"/txbeat/receive", //to put in tx pool
		TransactionBeatRecv,
	},
	Route{
		"TransactionPoolRecv", //api of bcHolder
		"GET",
		"/txbeat/allprev", //to put in tx pool
		TransactionPoolRecv,
	},
	Route{
		"GetMyId", //api of bcHolder
		"GET",
		"/getmyid", //to put in tx pool
		GetMyId,
	},
}
