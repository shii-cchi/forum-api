package app

import "log"
import "github.com/go-chi/chi"
import "github.com/shii-cchi/forum-api/internal/server"

func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()

	srv, err := server.NewServer(r)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
