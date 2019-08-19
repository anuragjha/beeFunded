package apps

type app interface {
	appServeReq()
	appServeRes()
}
