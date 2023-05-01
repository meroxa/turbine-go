package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

var exiter func(code int)

func init() {
	exiter = os.Exit

	buildCmd.StringVar(&buildGitSHA, "gitsha", "", "git commit sha used to reference the code deployed")
	buildCmd.StringVar(&buildListenAddr, "turbine-core-server", "", "address of the turbine core server")
	buildCmd.StringVar(&buildAppPath, "app-path", execPath(), "path to the turbine application")

	serverCmd.StringVar(&serverListenAddr, "serve-addr", os.Getenv("MEROXA_FUNCTION_ADDR"), "listen serve address")
	serverCmd.StringVar(&serverFuncName, "serve-func", "", "name of function to serve")
	serverCmd.SetOutput(os.Stderr)

	if len(os.Args) < 2 {
		usage()
	}
}

func parseFlags() string {
	cmd, args := os.Args[1], os.Args[2:]
	switch cmd {
	case serverCmdName:
		if err := serverCmd.Parse(args); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		serverCmd.VisitAll(requiredServerFlag)
	case buildCmdName:
		if err := buildCmd.Parse(args); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		buildCmd.VisitAll(requiredBuildFlag)
	default:
		usage()
	}

	return cmd
}

func execPath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("unable to locate executable path; error: %s\n", err)
	}
	return path.Dir(exePath)
}

func requiredServerFlag(f *flag.Flag) {
	if f.Value.String() == "" {
		fmt.Fprintf(os.Stderr, "error: -%s is required\n\n", f.Name)
		serverUsage()
		exiter(2)
	}
}

func requiredBuildFlag(f *flag.Flag) {
	if f.Value.String() == "" {
		fmt.Fprintf(os.Stderr, "error: -%s is required\n\n", f.Name)
		buildUsage()
		exiter(2)
	}
}

func buildUsage() {
	fmt.Fprintln(os.Stderr, "build:")
	buildCmd.PrintDefaults()
}

func serverUsage() {
	fmt.Fprintln(os.Stderr, "server:")
	serverCmd.PrintDefaults()
}

func usage() {
	buildUsage()
	serverUsage()
	exiter(2)
}
