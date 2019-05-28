package slack

import (
	"net/http"

	"github.com/hori-ryota/zaperr"
)

func (h HTTPHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error("handleError", zaperr.ToField(err))
	// TODO fix response
}
