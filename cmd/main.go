package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"webinar2609/internal/server"
	"webinar2609/internal/service/image"
	"webinar2609/internal/service/telegram"
)

func main() {
	imageService := image.NewImageService()
	telegramService := telegram.NewTelegramService("7937865545:AAHExRULKb9C1Y3sLF_kUogKYbtDjUoip1o", imageService)
	httpServer := server.NewServer(imageService)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Print("Telegram service starting...")
		telegramService.Start()
	}()

	go func() {
		httpServer.Start()
	}()
	<-done

	log.Print("Telegram service stopping...")
	telegramService.Stop()
	log.Print("Telegram service stopped.")

}
