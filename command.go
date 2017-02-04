package jobs

import (
	"context"
	"log"
	"os/exec"
	"time"
)

func CmdVisitable(cmdString string) Visitable {

	f := func(ctx context.Context) (context.Context, error) {

		ctxTimebound, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()

		cmd := exec.Command("bash", "-c", cmdString)
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		done := make(chan error)
		go func() { done <- cmd.Wait() }()
		select {
		case err := <-done:
			// exited
			return ctx, err
		case <-ctxTimebound.Done():
			// timed out
			return ctxTimebound, ctx.Err()
		}
	}

	return funcVisitable(f)
}
