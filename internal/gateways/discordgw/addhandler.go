package discordgw

import (
	"os"
	"os/signal"
	"syscall"

	"dstwilio/internal/core/ports"
)

// AddHandler adds a new message handler to the Discord gateway.
// It takes a `ports.NewMessageHandler` function as a parameter and returns an error.
// The function registers the handler with the Discord session and opens the connection.
// It then waits for a termination signal to gracefully close the connection.
// If an error occurs during the session opening or closing, it is logged, and the function returns the error.
func (d *Gateway) AddHandler(handler ports.NewMessageHandler) error {
	d.session.AddHandler(handler)

	if err := d.session.Open(); err != nil {
		d.logger.Error("error opening connection:", err, nil)

		return err
	}

	defer func() {
		if err := d.session.Close(); err != nil {
			d.logger.Error("error closing connection:", err, nil)

			os.Exit(1)
		}
	}()

	d.logger.Info("bot is now running. Press CTRL-C to exit.", nil)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	return nil
}
