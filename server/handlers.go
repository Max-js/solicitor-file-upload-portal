package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"strconv"
	"io"
)

type API struct {
	store          *Store
	maxUploadBytes int64
}

var allowedContentTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func (a *API) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/documents", a.createDocument)
	mux.HandleFunc("GET /api/documents", a.getDocuments)
	mux.HandleFunc("GET /api/client", a.getClient)
	return mux
}

func (a *API) createDocument(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, a.maxUploadBytes)
	if err := r.ParseMultipartForm(a.maxUploadBytes); err != nil {
		writeError(w, http.StatusBadRequest, "upload too large or malformed")
		return
	}

	clientID, err := strconv.ParseInt(r.FormValue("clientId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid or missing clientId")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()

	// INFO: Read the first 512 bytes and detect the file type from the bytes 
	head := make([]byte, 512)
	n, err := io.ReadFull(file, head)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		writeError(w, http.StatusBadRequest, "could not read file")
		return
	}

	contentType := http.DetectContentType(head[:n])
	if !allowedContentTypes[contentType] {
		writeError(w, http.StatusUnsupportedMediaType, "only PDF, JPEG, and PNG files are allowed")
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		writeError(w, http.StatusInternalServerError, "could not process file")
		return
	}

	storageKey := uuid.NewString()

	size, err := a.store.SaveFile(storageKey, file)
	if err != nil {
		log.Printf("save file: %v", err)
		writeError(w, http.StatusInternalServerError, "could not store file")
		return
	}

	doc := &Document{
		ClientID:    clientID,
		Filename:    header.Filename,
		ContentType: contentType,
		SizeBytes:   size,
		StorageKey:  storageKey,
		Status:      "pending",
	}

	if err := a.store.CreateDocument(r.Context(), doc); err != nil {
		log.Printf("create document: %v", err)
		writeError(w, http.StatusInternalServerError, "could not save document")
		return
	}

	writeJSON(w, http.StatusCreated, doc)
}

func (a *API) getDocuments(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("clientId")

	var (
		docs []Document
		err  error
	)

	if clientID != "" {
		docs, err = a.store.getDocumentsByClient(r.Context(), clientID)
	} else {
		docs, err = a.store.getAllDocuments(r.Context())
	}

	if err != nil {
		log.Printf("get documents: %v", err)
		writeError(w, http.StatusInternalServerError, "could not get documents")
		return
	}

	writeJSON(w, http.StatusOK, docs)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (a *API) getClient(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		writeError(w, http.StatusBadRequest, "missing email")
		return
	}

	client, err := a.store.GetClientByEmail(r.Context(), email)
	if err != nil {
		writeError(w, http.StatusNotFound, "client not found")
		return
	}
	
	writeJSON(w, http.StatusOK, client)
}