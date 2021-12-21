package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	dust "github.com/shellyln/dust-lang/cmd"
)

func printHelp() {
	name, _ := os.Executable()
	name = filepath.Base(name)
	fmt.Println("NAME:")
	fmt.Printf("  %s - demo app\n", name)
	fmt.Println()
	fmt.Println("   https://shellyln.github.io/")
	fmt.Println()
	fmt.Println("VERSION:")
	fmt.Println("  Version:", Version)
	fmt.Println("  Revision:", Revision)
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Printf("  %s [options] [subcommand] [subcmd-options] [files|arguments]\n", name)
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Printf("  %s run -e \"3 + 5\"\n", name)
	fmt.Printf("  %s run example.rs\n", name)
	fmt.Println()
	fmt.Println("OPTIONS:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("SUBCOMMANDS:")
	fmt.Println("  run   Run script")
	fmt.Println("  help  Show help")
}

func printRunHelp(runCmd *flag.FlagSet) {
	name, _ := os.Executable()
	name = filepath.Base(name)
	fmt.Println("USAGE:")
	fmt.Printf("  %s run [options] [files|arguments]\n", name)
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Printf("  %s run -e \"3 + 5\"\n", name)
	fmt.Printf("  %s run example.rs\n", name)
	fmt.Println()
	fmt.Println("OPTIONS:")
	runCmd.PrintDefaults()
}

func printVersion() {
	fmt.Println("Version:", Version)
	fmt.Println("Revision:", Revision)
}

func main() {
	var helpFlag bool
	var versionFlag bool
	var dustFlags dust.CommandFlags

	flag.BoolVar(&helpFlag, "h", false, "Show this help")
	flag.BoolVar(&helpFlag, "help", false, "Show this help")
	flag.BoolVar(&versionFlag, "v", false, "Print the version")
	flag.BoolVar(&versionFlag, "version", false, "Print the version")

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runCmd.BoolVar(&dustFlags.Eval, "e", false, "Evaluate the following argument as script")
	runCmd.BoolVar(&dustFlags.Eval, "eval", false, "Evaluate the following argument as script")

	flag.Parse()
	args := flag.Args()

	printSubHelp := func(subCmd string) {
		switch subCmd {
		case "run":
			printRunHelp(runCmd)
			os.Exit(0)
		case "help":
		case "version":
		}
	}

	if helpFlag {
		if len(args) > 0 {
			printSubHelp(args[0])
		}
		printHelp()
		os.Exit(0)
	}

	if versionFlag {
		printVersion()
		os.Exit(0)
	}

	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	} else {
		// subcommands
		switch args[0] {
		case "run":
			runCmd.Parse(os.Args[2:])
			dust.SubcommandMain(dustFlags, runCmd.Args())
		case "help":
			if len(args) > 1 {
				printSubHelp(args[1])
			}
			printHelp()
			os.Exit(0)
		case "version":
			printVersion()
			os.Exit(0)
		default:
			runCmd.Parse(os.Args[1:])
			dust.SubcommandMain(dustFlags, runCmd.Args())
		}
	}

	os.Exit(0)
}
