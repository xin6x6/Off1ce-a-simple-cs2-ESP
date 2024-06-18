package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ttacon/chalk"
)

func cliMenu() {
	for {
		//fmt.Print(chalk.Magenta.Color("          ____             \n  ___ ___|___ \\ ____  ___  \n / __/ __| __) / _  |/ _ \\ \n| (__\\__ \\/ __/ (_| | (_) |\n \\___|___/_____\\__, |\\___/ \n               |___/       \n"))
		fmt.Print(chalk.Red.Color("   _______                            _ \n"))
		fmt.Print(chalk.Red.Color("  |  _____|                          | | \n"))
		fmt.Print(chalk.Green.Color("  | |_____ _____  _____   ____  ____ | | \n"))
		fmt.Print(chalk.Green.Color("  |  _____/ ___/ |  __ \\ / ___ / __ \\| | \n"))
		fmt.Print(chalk.Cyan.Color("  | |_____\\____ \\|  ___/| (___| ____|| |__ \n"))
		fmt.Print(chalk.Cyan.Color("  |_______|_____/| |     \\____ \\_____|___/ \n"))
		fmt.Print(chalk.Magenta.Color("                 |_| \n"))

		fmt.Println(chalk.Cyan.Color("\t\t By Wh1nB --version 1.2\n"))
		if teamCheck {
			fmt.Println(chalk.Green.Color("[1] Team check [CT will show as blue box, T will show as red box] [ON]"))
		} else {
			fmt.Println(chalk.Red.Color("[1] Team check [CT will show as blue box, T will show as red box] [OFF]"))
		}
		if headCircle {
			fmt.Println(chalk.Green.Color("[2] Drawing head [Cyan] [ON]"))
		} else {
			fmt.Println(chalk.Red.Color("[2] Drawing head [Cyan] [OFF]"))
		}
		if skeletonRendering {
			fmt.Println(chalk.Green.Color("[3] Drawing skeleton [ON]"))
		} else {
			fmt.Println(chalk.Red.Color("[3] Drawing skeleton [OFF]"))
		}
		if boxRendering%3 == 1 { //corner
			fmt.Println(chalk.Green.Color("[4] Drawing box [Corner]"))
		}
		if boxRendering%3 == 2 { //full
			fmt.Println(chalk.Cyan.Color("[4] Drawing box [Full rectangle]"))
		}
		if boxRendering%3 == 0 { //off
			fmt.Println(chalk.Red.Color("[4] Drawing box [Off]"))
		}
		if healthBarRendering {
			fmt.Println(chalk.Green.Color("[5] Health bar [green] [ON]"))
		} else {
			fmt.Println(chalk.Red.Color("[5] Health bar [green] [OFF]"))
		}
		if healthTextRendering {
			fmt.Println(chalk.Green.Color("[6] Health number [ON]"))
		} else {
			fmt.Println(chalk.Red.Color("[6] Health number [OFF]"))
		}
		if nameRendering {
			fmt.Println(chalk.Green.Color("[7] Player Name [ON]"))
		} else {
			fmt.Println(chalk.Red.Color("[7] Player name [OFF]"))
		}
		if line {
			fmt.Println(chalk.Green.Color("[8] Drawing Line [ON]"))
		} else {
			fmt.Println(chalk.Red.Color("[8] Drawing Line [OFF]"))
		}
		fmt.Println(chalk.Cyan.Color("[9] Adjust frame delay [") + fmt.Sprint(frameDelay) + chalk.Cyan.Color("]"))
		fmt.Println(chalk.Red.Color("[10] Exit"))
		fmt.Print(chalk.Cyan.Color("[Enter your selection]: "))
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1":
			teamCheck = !teamCheck
		case "2":
			headCircle = !headCircle
		case "3":
			skeletonRendering = !skeletonRendering
		case "4":
			boxRendering++
		case "5":
			healthBarRendering = !healthBarRendering
		case "6":
			healthTextRendering = !healthTextRendering
		case "7":
			nameRendering = !nameRendering
		case "8":
			line = !line
		case "9":
			fmt.Println(chalk.Red.Color("Higer frame delay = lower performance impact but higher ESP latency"))
			fmt.Print(chalk.Cyan.Color("[Enter frame delay]: "))
			var delay uint32
			fmt.Scanln(&delay)
			frameDelay = delay
		case "10":
			os.Exit(0)
		default:
			fmt.Println(chalk.Red.Color("Invalid selection!"))
			time.Sleep(400 * time.Millisecond)
		}
		// Clear the console
		fmt.Print("\033[H\033[2J")
	}
}
