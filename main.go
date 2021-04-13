package main

import (
	"Netron1-Go/api"
	"Netron1-Go/config"
	"Netron1-Go/gui"
	"Netron1-Go/simulation"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const configFile = "config/config.json"

var mainLoop = true
var consoleLoop = true
var simLoop = true

func main() {
	config, err := config.NewConfig(configFile)

	if err != nil {
		log.Fatal(err)
	}

	// Our channels with the simulation coroutine
	chToSim := make(chan string)
	chFromSim := make(chan string)

	// -----------------------------------------------------
	// Setup GUI
	// -----------------------------------------------------
	surface := gui.NewSurfaceBuffer()

	surface.Open()
	surface.Configure()

	go surface.Run(chToSim, chFromSim)

	// -----------------------------------------------------
	// Setup console
	// -----------------------------------------------------
	fmt.Println("Welcome to Netron1 Go edition")
	fmt.Println("LogRoot: " + config.LogRoot())

	printHelp()

	go messageFromConsole(chToSim)

	go messageFromSim(chFromSim, config, surface)

	// -----------------------------------------------------
	// Setup simulation
	// -----------------------------------------------------
	sim := simulation.NewSimulation()

	sim.Initialize(surface.Raster(), surface)
	go sim.Start(chToSim, chFromSim)

	// -----------------------------------------------------
	// Stall until other coroutines finish
	// -----------------------------------------------------
	for mainLoop {
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Closing channels...")
	close(chFromSim)
	close(chToSim)

	fmt.Println("Closing window...")
	surface.Close()

	fmt.Println("Goodbye.")
}

func messageFromConsole(chToSim chan string) {
	reader := bufio.NewReader(os.Stdin)

	for consoleLoop {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		switch text {
		case "q":
			chToSim <- "exit"
			fmt.Println("-------")
			mainLoop = false
			consoleLoop = false
			simLoop = false
		case "r":
			chToSim <- "run"
		case "p":
			chToSim <- "pause"
		case "e":
			chToSim <- "step"
		case "u":
			chToSim <- "resume"
		case "s":
			chToSim <- "reset"
		case "t":
			chToSim <- "stop"
		case "a":
			chToSim <- "status"
		default:
			fmt.Println("*********************")
			fmt.Println("** Unknown command **")
			fmt.Println("*********************")
		}

		if !mainLoop {
			printHelp()
		}
	}
}

func messageFromSim(chFromSim chan string, config api.IConfig, surface api.ISurface) {

	for simLoop {
		select {
		case msg := <-chFromSim:
			switch msg {
			case "Exited":
				fmt.Println("Simulation exited.")
				config.SetExitState("Exited")
				simLoop = false
				mainLoop = false
			case "Started":
				fmt.Println("Simulation started.")
			case "Stepped":
				fmt.Println("Simulation stepped.")
			case "Reset":
				fmt.Println("Simulation reset.")
			case "Terminated":
				fmt.Println("Simulation terminated.")
				config.SetExitState("Terminated")
			case "Paused":
				fmt.Println("Simulation paused.")
				config.SetExitState("Paused")
			case "Stopped":
				fmt.Println("Simulation stopped.")
				config.SetExitState("Stopped")
			case "Complete":
				fmt.Println("Simulation completed.")
				config.SetExitState("Completed")
			case "Running":
				fmt.Println("Simulation is ready running!")
			case "Not Running":
				fmt.Println("Simulation isn't running.")
			case "Already Paused":
				fmt.Println("Simulation is already paused.")
			case "Not Paused":
				fmt.Println("Simulation is not paused.")
			case "Resumed":
				fmt.Println("Simulation resumed.")
			default:
				fmt.Println(msg)
			}
		}

		if simLoop {
			fmt.Print("> ")
		}
	}
}

func printHelp() {
	fmt.Println("-----------------------------")
	fmt.Println("Commands:")
	fmt.Println("  q: quit")
	fmt.Println("  r: run simulation")
	fmt.Println("  p: pause simulation")
	fmt.Println("  u: resume simulation")
	fmt.Println("  s: reset simulation.")
	fmt.Println("  t: stop simulation")
	fmt.Println("  a: status of simulation")
	fmt.Println("  h: this help menu")
	fmt.Println("-----------------------------")
	fmt.Print("> ")
}
