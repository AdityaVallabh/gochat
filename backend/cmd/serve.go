/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AdityaVallabh/gochat/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the chat server and serve requests",
	Run: func(cmd *cobra.Command, args []string) {
		if err := serve(viper.GetString("ADDRESS")); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(addr string) error {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"), viper.GetInt("DB_PORT"), "disable", "Asia/Kolkata")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	s := server.Server{
		Router: mux.NewRouter(),
		DB:     db,
	}
	if err := s.Setup(); err != nil {
		return err
	}
	http.HandleFunc("/", s.ServeHTTP)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Info("now listening on ", l.Addr())
	err = http.Serve(l, nil)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	log.Info("Stopped server")
	return nil
}
