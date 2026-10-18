/*
 * AF: AKMA Application Function of the MVP (3GPP TS 33.535), serving the
 * Ua* protocol of docs/af-side/ua-star.md.
 */

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/florianoneto07/akma-free5gc/af/af/internal/aanf"
	"github.com/florianoneto07/akma-free5gc/af/af/internal/config"
	"github.com/florianoneto07/akma-free5gc/af/af/internal/logger"
	"github.com/florianoneto07/akma-free5gc/af/af/internal/server"
	logger_util "github.com/free5gc/util/logger"
)

func main() {
	cfgPath := flag.String("c", "config/afcfg.yaml", "configuration `file`")
	logPath := flag.String("l", "", "also write the log to `file`")
	flag.Parse()
	if err := run(*cfgPath, *logPath); err != nil {
		logger.MainLog.Errorf("AF: %v", err)
		os.Exit(1)
	}
}

func run(cfgPath, logPath string) error {
	if logPath != "" {
		if err := logger_util.LogFileHook(logger.Log, logPath); err != nil {
			return err
		}
	}
	s, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("configuration %s: %w", cfgPath, err)
	}
	logger.CfgLog.Infof("Read config from [%s]", cfgPath)
	if err := setLogger(s.Logger); err != nil {
		return err
	}
	logger.MainLog.Infof("AF_ID %s; realm %q; qop %v; nonce lifetime %s; K_AF expiry policy %s; AAnF %s",
		s.AfID.Transport(), s.Realm, s.Qop, s.NonceLifetime, s.KafExpiry, s.AanfURI)

	srv := &http.Server{
		Addr:              s.Listen,
		Handler:           server.New(s, aanf.New(s.AanfURI, s.AfID)),
		ReadHeaderTimeout: 10 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()

	logger.MainLog.Infof("Ua* server listening on %s for http://%s%s", s.Listen, s.FQDN, s.Path)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	logger.MainLog.Infof("AF terminated")
	return nil
}

func setLogger(l config.Logger) error {
	if !l.Enable {
		logger.Log.SetOutput(nopWriter{})
	}
	if l.Level != "" {
		lvl, err := logrus.ParseLevel(l.Level)
		if err != nil {
			return fmt.Errorf("logger.level: %w", err)
		}
		logger.Log.SetLevel(lvl)
	}
	logger.Log.SetReportCaller(l.ReportCaller)
	return nil
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }
