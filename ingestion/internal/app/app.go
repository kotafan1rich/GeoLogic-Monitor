package app

type App struct {
	di *diContainer
}

func New() *App {
	return &App{di: newDIContainer()}
}

// func (a *App) Run(ctx context.Context) error {
// 	stopCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
// 	defer stop()

// 	<-stopCtx.Done()

// 	stop()

// 	closerCtx, closerCancel := context.WithTimeout(ctx, closeDuration)
// 	defer closerCancel()

// 	if err := closer.CloseAll(closerCtx); err != nil {
//
// 	}

// 	return nil
// }
