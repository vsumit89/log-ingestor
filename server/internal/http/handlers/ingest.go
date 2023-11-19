package handlers

import (
	"encoding/json"
	"logswift/internal/dtos"
	"logswift/pkg/helper"
	"net/http"
)

func (h *logHandler) ingest(w http.ResponseWriter, r *http.Request) {
	var logEntry dtos.LogEntry
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	h.logSvc.IngestLog(logEntry)
	// if err != nil {
	// 	helper.SendJSONResponse(w, http.StatusInternalServerError, "failed to ingest log", nil)
	// 	return
	// }

	helper.SendJSONResponse(w, http.StatusOK, "log ingested successfully", nil)
}
