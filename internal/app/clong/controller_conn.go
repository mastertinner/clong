package clong

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// ControllerConnHandler handles a WebSocket connection coming from a controller.
func ControllerConnHandler(h *Hub, up websocket.Upgrader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, errUpgradingConnection))
			return
		}
		defer func() {
			err := ws.Close()
			if err != nil {
				log.Fatal(errors.Wrap(err, "error closing websocket"))
			}
		}()

		h.RegisterController <- ws

		for {
			var c Control
			id, ok := cookieVal(r.Cookies(), "userid")
			if !ok {
				handleHTTPError(w, ErrUserIDMissing)
				h.UnregisterController <- ws
				break
			}
			name, ok := cookieVal(r.Cookies(), "username")
			if !ok {
				handleHTTPError(w, ErrUserNameMissing)
				h.UnregisterController <- ws
				break
			}
			c.Player = User{
				ID:   id,
				Name: name,
			}

			err = ws.ReadJSON(&c)
			if err != nil {
				handleHTTPError(w, errors.Wrap(err, errReadingJSON))
				h.UnregisterController <- ws
				break
			}

			h.Controls <- c
		}
	})
}

// cookieVal returns the value of a cookie.
func cookieVal(cookies []*http.Cookie, name string) (string, bool) {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value, true
		}
	}

	return "", false
}