package lifecycle

var (
	// OnConnect is the life cycle phase for connecting to stateful services.
	OnConnect = New("connect")

	// OnStart is the life cycle phase for service start-up.
	OnStart = New("start")

	// OnShutdown is the life cycle phase for service shutdown.
	OnShutdown = New("shutdown")
)
