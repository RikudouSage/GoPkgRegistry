package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go.chrastecky.dev/go-pkg-repository/db"
	"go.chrastecky.dev/go-pkg-repository/dto"
	"go.chrastecky.dev/go-pkg-repository/helper"
)

func StorePackageHandler(writer http.ResponseWriter, req *http.Request, database *db.Client) {
	pkg := &dto.Package{}
	if err := json.NewDecoder(req.Body).Decode(&pkg); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		helper.WriteJSON(writer, map[string]string{
			"error": "internal error",
		})
		return
	}

	if !pkg.IsValid() {
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteJSON(writer, map[string]string{
			"error": "invalid package, some required fields are missing",
		})
		return
	}

	if err := database.StorePackage(pkg); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		helper.WriteJSON(writer, map[string]string{
			"error": "internal error",
		})
		return
	}

	writer.WriteHeader(http.StatusCreated)
	helper.WriteJSON(writer, pkg)
}
