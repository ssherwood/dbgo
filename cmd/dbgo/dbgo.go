package main

import (
	"context"
	"dbgo/internal/common"
	"dbgo/internal/generator"
	"dbgo/internal/providers"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"log"
	"math/rand"
	"os"
	"time"
)

type ProgramArgs struct {
	ConfigFile string
	Iterations int
	Mode       string
	Provider   string
}

func main() {
	fmt.Println("DB Go!")

	programArgs := initArgs()
	dbConfig := loadDBConfig(programArgs.ConfigFile)
	dbProvider := getDatabaseProvider(programArgs, dbConfig)
	if err := dbProvider.Init(); err != nil {
		fmt.Println("Unable to initialize provider:", err)
		os.Exit(1)
	}
	defer dbProvider.Close()

	if programArgs.Mode == "discover" {
		fmt.Println("DB Discover Mode: " + dbConfig.Schema)

		tableConfigs, err := dbProvider.Discover(dbConfig.Schema)
		if err != nil {
			fmt.Println("Unable to discover the database metadata", err)
			os.Exit(1)
		}

		if yamlData, err := yaml.Marshal(&tableConfigs); err != nil {
			fmt.Println("Unable to parse the database metadata into yaml", err)
			os.Exit(1)
		} else {
			fmt.Println(string(yamlData))
		}
	} else {
		fmt.Printf("DB Init Mode: %s://%s:%d/%s?schema=%s\n", dbConfig.Name, dbConfig.Host, dbConfig.Port, "TODO", dbConfig.Schema)

		for tableIdx := 0; tableIdx < len(dbConfig.Tables); tableIdx++ {
			start := time.Now()
			columns := generator.ProcessColumns(dbConfig.Tables[tableIdx].Columns)

			//runBatchInserts(dbConfig, dbProvider, dbConfig.Tables[tableIdx].Name, columns, programArgs.Iterations, 1000)
			runSingleRowInserts(dbConfig, dbProvider, dbConfig.Tables[tableIdx].Name, columns, programArgs.Iterations)

			elapsed := time.Since(start)
			fmt.Printf("Inserted %d row in %v.\n", programArgs.Iterations, elapsed)
		}
	}
}

// runSingleRowInserts
//
// This implementation uses a channel that will block until a specific number of messages
// are received.  This effectively throttles the number of concurrent go routines to keep
// from over saturating the CPU and providers connections.
func runSingleRowInserts(dbConfig common.DatabaseConfig, dbProvider providers.DatabaseProvider, tableName string, columns []generator.ColumnDefinition, iterations int) {
	errGrp, ctx := errgroup.WithContext(context.Background())
	blockingChannel := make(chan struct{}, 128)

	// TODO should check to see if the provider supports parallel operations...

	for k := 0; k < iterations; k++ {
		blockingChannel <- struct{}{} // block here after each 128 goroutine calls

		row := k // see https://www.calhoun.io/gotchas-and-common-mistakes-with-closures-in-go/#variables-declared-in-for-loops-are-passed-by-reference

		errGrp.Go(func() error {
			select {
			case <-ctx.Done():
				<-blockingChannel
				return ctx.Err()
			default:
				err := dbProvider.InsertSingleRow(dbConfig.Schema, tableName, columns, row, false)
				<-blockingChannel
				return err
			}
		})
	}

	fmt.Println("Finalizing...")

	if err := errGrp.Wait(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// runBatchInserts
func runBatchInserts(dbConfig common.DatabaseConfig, dbProvider providers.DatabaseProvider, tableName string, columns []generator.ColumnDefinition, iterations int, batchSize int) {
	errGrp, ctx := errgroup.WithContext(context.Background())
	blockingChannel := make(chan struct{}, 32)

	if batchSize > iterations {
		batchSize = iterations
	}

	for row := 0; row < iterations; row += batchSize {
		blockingChannel <- struct{}{} // block until # or messages received

		startingRow := row
		if startingRow+batchSize > iterations {
			batchSize = iterations - startingRow
		}

		errGrp.Go(func() error {
			select {
			case <-ctx.Done():
				<-blockingChannel
				return ctx.Err()
			default:
				err := dbProvider.InsertBatchRow(dbConfig.Schema, tableName, columns, startingRow, batchSize)
				<-blockingChannel
				return err
			}
		})
	}

	fmt.Println("Finalizing...")

	if err := errGrp.Wait(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func initArgs() ProgramArgs {
	modePtr := flag.String("mode", "discover", "mode of operation {discover, init}")
	providerPtr := flag.String("provider", "stdio", "target provider {stdio, csv, ysql}")
	configFilePtr := flag.String("config", "sample.yml", "the yaml config file")
	iterationsPtr := flag.Int("iters", 1, "number of iterations to run")
	flag.Parse()

	return ProgramArgs{
		ConfigFile: *configFilePtr,
		Iterations: *iterationsPtr,
		Mode:       *modePtr,
		Provider:   *providerPtr,
	}
}

func loadDBConfig(configFile string) common.DatabaseConfig {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var dbConfig common.DatabaseConfig
	if err = dbConfig.Parse(data); err != nil {
		log.Fatal(err)
	}

	if dbConfig.Seed != 0 {
		rand.New(rand.NewSource(dbConfig.Seed))
	}

	return dbConfig
}

func getDatabaseProvider(programArgs ProgramArgs, dbConfig common.DatabaseConfig) providers.DatabaseProvider {
	switch programArgs.Provider {
	case "ysql":
		return &providers.YugabyteYSQL{Config: dbConfig}
	default:
		return &providers.StdoutDBProvider{Writer: os.Stdout}
	}
}
