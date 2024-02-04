package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	instantaneousPowerGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ipmi_instantaneous_power_watts",
		Help: "Instantaneous power reading in watts",
	})
	minimumPowerGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ipmi_minimum_power_watts",
		Help: "Minimum power reading in watts during the sampling period",
	})
	maximumPowerGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ipmi_maximum_power_watts",
		Help: "Maximum power reading in watts during the sampling period",
	})
	averagePowerGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ipmi_average_power_watts",
		Help: "Average power reading in watts over the sample period",
	})
)

func init() {
	// Register the Prometheus gauges
	prometheus.MustRegister(instantaneousPowerGauge)
	prometheus.MustRegister(minimumPowerGauge)
	prometheus.MustRegister(maximumPowerGauge)
	prometheus.MustRegister(averagePowerGauge)

	// Initialize Viper
	viper.SetConfigName("ipmi_power_exporter")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath(".")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Error reading config file")
	}

	// Set up logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := viper.GetString("log.level")
	logFormat := viper.GetString("log.format")

	switch logFormat {
	case "json":
		log.Logger = log.Output(os.Stderr)
	case "pretty":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	default:
		log.Fatal().Msg("Invalid log format specified")
	}

	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		log.Fatal().Msg("Invalid log level specified")
	}
}

func collectPowerMetrics(ctx context.Context, host string, user string, pass string) {
	// Regular expressions to match the power readings
	instantaneousRegex := regexp.MustCompile(`Instantaneous power reading:\s+(\d+) Watts`)
	minimumRegex := regexp.MustCompile(`Minimum during sampling period:\s+(\d+) Watts`)
	maximumRegex := regexp.MustCompile(`Maximum during sampling period:\s+(\d+) Watts`)
	averageRegex := regexp.MustCompile(`Average power reading over sample period:\s+(\d+) Watts`)

	// Extract the power readings
	extractValue := func(regex *regexp.Regexp, metric prometheus.Gauge, outputStr string) float64 {
		match := regex.FindStringSubmatch(outputStr)
		if len(match) == 2 {
			value, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				log.Error().Err(err).Msg("Error parsing power reading")

				return -1
			}
			metric.Set(value)

			return value
		}

		return -1
	}

	// Building command line params
	params := []string{}
	if host != "" {
		params = []string{"-I", "lanplus", "-H", host, "-U", user, "-P", pass}
	}

	params = append(params, "dcmi", "power", "reading")

	log.Info().Msg("Collecting data")
	log.Debug().Interface("params", params).Msg("calling ipmitools with params")

	// Execute the ipmitool command with user and password
	cmd := exec.CommandContext(ctx, "ipmitool", params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Error executing ipmitool")
		return
	}

	// Convert the output to a string
	outputStr := string(output)

	// Extract and set the power readings
	cval := extractValue(instantaneousRegex, instantaneousPowerGauge, outputStr)
	minval := extractValue(minimumRegex, minimumPowerGauge, outputStr)
	maxval := extractValue(maximumRegex, maximumPowerGauge, outputStr)
	avgval := extractValue(averageRegex, averagePowerGauge, outputStr)

	log.Debug().
		Str("current", fmt.Sprintf("current: %0.2f Watt", cval)).
		Str("min", fmt.Sprintf("min: %0.2f Watt", minval)).
		Str("max", fmt.Sprintf("max: %0.2f Watt", maxval)).
		Str("avg", fmt.Sprintf("avg: %0.2f Watt", avgval)).
		Msg("Collected power metrics")
}
func main() {

	// Get the IPMI configuration
	host := viper.GetString("ipmi.host")
	user := viper.GetString("ipmi.user")
	pass := viper.GetString("ipmi.pass")

	// Get the collect interval and listen port from the configuration
	collectInterval := viper.GetDuration("collect.interval")
	listenPort := viper.GetInt("server.port")
	listenAddress := viper.GetString("server.address")

	// Start polling in the background
	go func() {
		ticker := time.NewTicker(collectInterval)
		defer ticker.Stop()

		// Create a context that can be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Set up a signal handler to cancel the context
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		// Start a goroutine to periodically collect metrics
		for {
			select {
			case <-ticker.C:
				collectPowerMetrics(ctx, host, user, pass)
			case <-ctx.Done():
				os.Exit(-1)
			case sig := <-sigCh:
				log.Info().Str("signal", sig.String()).Msg("Received signal, shutting down")
				os.Exit(0)
			}
		}
	}()

	// Expose the registered metrics via HTTP
	addr := fmt.Sprintf("%s:%d", listenAddress, listenPort)
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msgf("starting power consume exporter daemon on %s", addr)
	log.Fatal().Err(http.ListenAndServe(addr, nil)).Msg("power consume exporter daemon stopped")
}
