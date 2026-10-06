package wecom

import (
	"net/http"

	"github.com/go-sphere/wecom-aibot-go-sdk/aibot"
)

func Handle(client *aibot.WSClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
