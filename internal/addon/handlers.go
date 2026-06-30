package addon

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func HandleCatalog(instance *Instance) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("path")
		parts := strings.Split(strings.Trim(path, "/"), "/")

		if len(parts) < 3 || parts[0] != "catalog" {
			http.Error(w, "Invalid catalog path", http.StatusBadRequest)
			return
		}

		catalogType := parts[1]
		catalogID := parts[2]

		extra := make(map[string]string)
		if len(parts) > 3 {
			extraFile := strings.TrimSuffix(parts[3], ".json")
			parseExtraParams(extraFile, extra)
		}

		query := r.URL.Query()
		if v := query.Get("genre"); v != "" {
			extra["genre"] = v
		}
		if v := query.Get("search"); v != "" {
			extra["search"] = v
		}
		if v := query.Get("skip"); v != "" {
			extra["skip"] = v
		}

		metas := instance.GetCatalog(catalogType, catalogID, extra)
		if metas == nil {
			metas = []MetaPreview{}
		}

		resp := CatalogResponse{Metas: metas}
		writeJSON(w, resp)
	}
}

func HandleStream(instance *Instance) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("path")
		parts := strings.Split(strings.Trim(path, "/"), "/")

		if len(parts) < 3 || parts[0] != "stream" {
			http.Error(w, "Invalid stream path", http.StatusBadRequest)
			return
		}

		streamType := parts[1]
		itemID := strings.TrimSuffix(parts[2], ".json")

		var stream *Stream
		if strings.HasPrefix(itemID, "iptv_series_ep_") {
			stream = instance.GetEpisodeStream(itemID)
		} else {
			stream = instance.GetStream(streamType, itemID)
		}

		var streams []Stream
		if stream != nil {
			streams = append(streams, *stream)
		}

		resp := StreamResponse{Streams: streams}
		writeJSON(w, resp)
	}
}

func HandleMeta(instance *Instance) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("path")
		parts := strings.Split(strings.Trim(path, "/"), "/")

		if len(parts) < 3 || parts[0] != "meta" {
			http.Error(w, "Invalid meta path", http.StatusBadRequest)
			return
		}

		metaType := parts[1]
		itemID := strings.TrimSuffix(parts[2], ".json")

		meta := instance.GetMeta(metaType, itemID)

		resp := MetaResponse{Meta: meta}
		writeJSON(w, resp)
	}
}

func HandleManifest(instance *Instance) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		manifest := instance.GetManifest()
		writeJSON(w, manifest)
	}
}

func parseExtraParams(extraFile string, extra map[string]string) {
	extraFile = strings.TrimSuffix(extraFile, ".json")
	if extraFile == "" {
		return
	}

	parts := strings.Split(extraFile, "=")
	if len(parts) == 2 {
		key := parts[0]
		value := parts[1]

		switch key {
		case "genre", "search", "skip":
			extra[key] = value
		}
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(v)
}

func ParseSkipParam(s string) int {
	if s == "" {
		return 0
	}
	skip, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return skip
}
