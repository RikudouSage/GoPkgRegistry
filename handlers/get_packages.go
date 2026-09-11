package handlers

import (
	"log"
	"net/http"

	"go.chrastecky.dev/go-pkg-repository/db"
	"go.chrastecky.dev/go-pkg-repository/helper"
)

// GetPackagesHandler writes all stored packages as JSON.
func GetPackagesHandler(writer http.ResponseWriter, _ *http.Request, database *db.Client) {
	pkgs, err := database.GetPackages()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		helper.WriteJSON(writer, map[string]string{
			"error": "internal error",
		})
		return
	}

	helper.WriteJSON(writer, pkgs)
}
