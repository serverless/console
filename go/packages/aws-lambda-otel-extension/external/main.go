package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"aws-lambda-otel-extension/external/extension"
	"aws-lambda-otel-extension/external/lib"
	"aws-lambda-otel-extension/external/logs"
	"aws-lambda-otel-extension/external/metrics"
	"aws-lambda-otel-extension/external/reporter"

	"go.uber.org/zap"
)

// INITIAL_QUEUE_SIZE is the initial size set for the synchronous logQueue
const INITIAL_QUEUE_SIZE = 100

var (
	extensionClient = extension.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
)

func main() {
	startTime := time.Now()

	ctx, cancel := context.WithCancel(context.Background())

	extWG := new(sync.WaitGroup)
	startWG := new(sync.WaitGroup)
	logger := lib.NewLogger()
	defer logger.Sync()

	// register extension while we start other services
	extWG.Add(1)
	go func() {
		// Register the extension
		err := extensionClient.Register(ctx)
		if err != nil {
			panic(err)
		}
		extWG.Done()
	}()

	settings, err := lib.GetExtensionSettings()
	if err != nil {
		logger.Error("Failed to get user settings", zap.Error(err))
		return
	}

	reportAgent := reporter.NewReporterClient(&settings)
	currentRequestData := reporter.NewCurrentRequestData(reportAgent)

	// Create Logs API agent
	logsApiAgent, err := logs.NewLogsApiAgent(reportAgent, currentRequestData, &settings)
	if err != nil {
		logger.Fatal("couldnt create logs api agent", zap.Error(err))
	}

	// Track OS signals to gracefully shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		logger.Debug(fmt.Sprintf("Received signal: %s, Exiting", s))
	}()

	// Create metrics agent/listener
	metricsApiListener := metrics.NewInternalHttpListener(reportAgent, currentRequestData)
	metricsApiListener.Start()

	// Subscribe to logs API
	// Logs start being delivered only after the subscription happens.
	startWG.Add(1)
	go func() {
		extWG.Wait()
		err = logsApiAgent.Init(extensionClient.ExtensionID)
		if err != nil {
			logger.Fatal("couldnt init logs api agent", zap.Error(err))
		}
		startWG.Done()
	}()

	startWG.Wait()

	// Will block until shutdown event is received or cancelled via the context.
	reportAgent.ReportInitDuration(startTime)

	// first call init next/event
	processEvents(ctx, logger, reportAgent, currentRequestData)

	reportAgent.RegisterClockStart()
	endWG := new(sync.WaitGroup)

	logger.Debug("Exiting")
	endWG.Add(1)
	go func() {
		logsApiAgent.Shutdown(ctx)
		endWG.Done()
	}()

	endWG.Add(1)
	go func() {
		metricsApiListener.Shutdown()
		endWG.Done()
	}()

	endWG.Add(1)
	go func() {
		err := reportAgent.Shutdown()
		if err != nil {
			logger.Error("Failed to shutdown report agent", zap.Error(err))
		}
		endWG.Done()
	}()

	endWG.Wait()

	reportAgent.ReportShutdownDuration()

	return
}

func processEvents(ctx context.Context, logger *lib.Logger, reportAgent *reporter.ReporterClient, currentRequestData *reporter.CurrentRequestData) {

	next := func() error {
		res, err := extensionClient.NextEvent(ctx)
		if err != nil {
			logger.Error(fmt.Sprintf("Error event: %s, Exiting", err))
			return err
		}
		// Exit if we receive a SHUTDOWN event
		if res.EventType == extension.Shutdown {
			logger.Debug("Received SHUTDOWN event, Exiting")
			return errors.New("SHUTDOWN")
		} else if res.EventType == extension.Invoke {
			currentRequestData.SetUniqueName("invoke")
			logger.Debug("Invoke received")
		}
		return nil
	}

	// first next to flag we're ready
	if err := next(); err != nil {
		return
	}

	for {
		// block until we receive runtimeDone
		reportAgent.WaitDone()
		reportAgent.WaitRequests(time.Millisecond * 1000)
		select {
		case <-ctx.Done():
			logger.Debug("Context cancelled, exiting")
			return
		default:
			logger.Debug("Reading next event...")
			reportAgent.ReportOverheadDuration()
			if err := next(); err != nil {
				return
			}
		}
	}
}
