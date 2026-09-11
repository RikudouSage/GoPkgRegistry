package handlers

import (
	"GoPkgRepository/db"
	"GoPkgRepository/helper"
	"log"
	"net/http"
)

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
