package handlers

import (
	"net/http"
	"strings"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
)

func DisconnectPlatformHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := getAuthenticatedUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	platform := strings.TrimSpace(r.URL.Query().Get("platform"))
	switch platform {
	case "x":
		if err := db.ClearUserTwitterInfo(user.ID); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to disconnect X")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "X disconnected"})
	case "linkedin":
		if err := db.ClearUserLinkedInInfo(user.ID); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to disconnect LinkedIn")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "LinkedIn disconnected"})
	default:
		writeError(w, http.StatusBadRequest, "Unsupported platform")
	}
}
