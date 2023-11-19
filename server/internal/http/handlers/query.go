package handlers

import (
	"logswift/internal/dtos"
	"logswift/pkg/helper"
	"logswift/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

func (h *logHandler) query(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetInstance()

	var err error
	// fetching filters and query from the request query params
	values := r.URL.Query()

	query := values.Get("q")

	var filters dtos.LogQueryFilters

	if values.Get("level") != "" {
		level := values.Get("level")
		filters.Level = &level
	}

	if values.Get("message") != "" {
		message := values.Get("message")
		filters.Message = &message
	}

	if values.Get("resourceId") != "" {
		resourceId := values.Get("resourceId")
		filters.ResourceID = &resourceId
	}

	if values.Get("traceId") != "" {
		traceId := values.Get("traceId")
		filters.TraceID = &traceId
	}

	if values.Get("spanId") != "" {
		spanId := values.Get("spanId")
		filters.SpanID = &spanId
	}

	if values.Get("commit") != "" {
		commit := values.Get("commit")
		filters.Commit = &commit
	}

	if values.Get("metadata_parent_id") != "" {
		parentResourceID := values.Get("metadata_parent_id")
		filters.ParentResourceID = &parentResourceID
	}

	if values.Get("from") != "" {
		parseFrom, err := time.Parse(time.RFC3339, values.Get("from"))
		if err != nil {
			logger.Error("Failed to parse time", "error", err)
			helper.SendJSONResponse(w, http.StatusBadRequest, "invalid time format", nil)
			return
		}
		filters.From = &parseFrom
	}

	if values.Get("to") != "" {
		parseTo, err := time.Parse(time.RFC3339, values.Get("to"))
		if err != nil {
			logger.Error("Failed to parse time", "error", err)
			helper.SendJSONResponse(w, http.StatusBadRequest, "invalid time format", nil)
			return
		}
		filters.To = &parseTo
	}

	if values.Get("limit") != "" {
		limit, err := strconv.Atoi(values.Get("limit"))
		if err != nil {
			logger.Error("Failed to parse limit", "error", err)
			helper.SendJSONResponse(w, http.StatusBadRequest, "invalid limit", nil)
			return
		}
		filters.Limit = limit
	} else {
		filters.Limit = 20
	}

	if values.Get("page") != "" {
		page, err := strconv.Atoi(values.Get("page"))
		if err != nil {
			logger.Error("Failed to parse page", "error", err)
			helper.SendJSONResponse(w, http.StatusBadRequest, "invalid page", nil)
			return
		}
		filters.Page = page
	} else {
		filters.Page = 1
	}

	// fetching logs from the storage
	result, err := h.logSvc.QueryLogs(query, filters)
	if err != nil {
		logger.Error("Failed to query logs", "error", err)
		helper.SendJSONResponse(w, http.StatusInternalServerError, "failed to query logs", nil)
		return
	}

	// writing the response
	helper.SendJSONResponse(w, http.StatusOK, "logs fetched successfully", result)
}
