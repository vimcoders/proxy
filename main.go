package main

import (
	"context"
	"os"
	"os/signal"
	"proxy/handler"
	"runtime"
	"syscall"

	"github.com/vimcoders/go-driver/log"
)

func main() {
	log.Info("NumCPU: ", runtime.NumCPU())
	quit := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	handler := handler.MakeHandler(ctx)
	go handler.ListenAndServe(ctx)
	log.Info("proxy running")
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
		log.Info("shutdown ->", s.String())
		cancel()
		handler.Close()
	default:
		log.Info("os.Signal ->", s.String())
	}
}
