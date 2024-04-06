package music_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

const musicDirectory = "./music"

func Run() {
	logger := httplog.NewLogger("music-stream-api", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		MessageFieldName: "msg",
	})

	logger.Info("Stream music API running...")

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/stream/{songTitle}", streamSong)

	r.Route("/songs", func(r chi.Router) {
		r.Post("/", uploadSong)
		r.Get("/", listSongs)
	})

	http.ListenAndServe(":3000", r)
}

func streamSong(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())

	urlParts := strings.Split(r.URL.Path, "/")
	songTitle := urlParts[len(urlParts)-1]

	musicFilePath := musicDirectory + "/" + songTitle

	musicFile, err := os.Open(musicFilePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer musicFile.Close()

	w.Header().Set("Content-Type", "audio/mpeg")

	chunkSize := 1024 * 512 // 0.5 MB

	buffer := make([]byte, chunkSize)

	for {
		bytesRead, err := musicFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				oplog.Info("All data have been read")
				break
			} else {
				http.Error(w, "Failed to read music file", http.StatusInternalServerError)
				return
			}
		}

		_, err = w.Write(buffer[:bytesRead])
		if err != nil {
			break
		}

		// Flush the response writer to ensure chunk is sent immediately
		w.(http.Flusher).Flush()

		msg := fmt.Sprintf("Sent chunk of size %d bytes\n", bytesRead)
		oplog.Info(msg)
	}

	oplog.Info("All data have been sent")
}

func uploadSong(w http.ResponseWriter, r *http.Request) {
	var in SongRequest

	err := json.NewDecoder(r.Body).Decode(&in)

	prettyJson, _ := json.MarshalIndent(&in, "", "  ")
	fmt.Println(string(prettyJson))

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(SongResponse{"title"})
	if err != nil {
		fmt.Printf("failed to encode created note: %v", err)
		return
	}
}

func listSongs(w http.ResponseWriter, _ *http.Request) {
	files, err := os.ReadDir(musicDirectory)
	if err != nil {
		fmt.Printf("failed to read files: %v", err)
	}

	var songs []SongResponse

	for _, file := range files {
		songs = append(songs, SongResponse{Title: file.Name()})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(songs)
	if err != nil {
		fmt.Printf("failed to encode created note: %v", err)
		return
	}
}

type SongRequest struct {
	Title string `json:"title"`
}

type SongResponse struct {
	Title string `json:"title"`
}
