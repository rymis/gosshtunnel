package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Get default value of variable from environment and create flag.
func envString(name, env, help string) *string {
	val := os.Getenv(env)
	set := ""
	if val != "" {
		set = " [set]"
	}
	helpMessage := fmt.Sprintf("%s (env: %s)%s", help, env, set)

	return flag.String(name, val, helpMessage)
}

type redirectInfo struct {
	Listen, Connect string
}

// Parse redirection
func newRedirectInfo(r string) (*redirectInfo, error) {
	parts := strings.Split(r, ":")
	res := &redirectInfo{}

	if len(parts) == 1 || len(parts) > 4 {
		return nil, fmt.Errorf("Waiting for redirect info in a form: [remoteHost:]remotePort:[localHost:]localPort")
	}

	if len(parts) == 3 {
		return nil, fmt.Errorf("Both remoteHost and localHost must be either present or absent")
	}

	if len(parts) == 2 { // "port:port"
		res.Listen = fmt.Sprintf("localhost:%s", parts[0])
		res.Listen = fmt.Sprintf("localhost:%s", parts[1])
	} else if len(parts) == 4 {
		res.Listen = fmt.Sprintf("%s:%s", parts[0], parts[1])
		res.Listen = fmt.Sprintf("%s:%s", parts[2], parts[3])
	}

	return res, nil
}

// Main function
func mainErr() error {
	key := envString("key", "SSH_PRIVATE_KEY", "Private key to use")
	host := envString("host", "SSH_REMOTE_HOST", "Remote host to connect. Can contain port (host:port).")
	user := envString("user", "USER", "Username to use for connection")
	blueprint := envString("blueprint", "SSH_KEY_BLUEPRINT", "Blueprint of server public key")
	genBlueprint := flag.Bool("showBlueprint", false, "Show remote server key blueprint and exit")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: gosshtunnel -key private_key -host remote-host.net [-user user] [-blueprint blueprint] [-showBlueprint] redirects...\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Each redirect has form [remoteHost:]remotePort:[localHost:]localhost\n")
		fmt.Fprintf(flag.CommandLine.Output(), "remoteHost and localHost must be either both present and both absent.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Additional way to define redirects is to use SSH_REDIRECTS environment variable.")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *key == "" {
		return fmt.Errorf("Argument `key` is mandatory")
	}

	if *host == "" {
		return fmt.Errorf("Argument `host` is mandatory")
	}

	if *user == "" {
		return fmt.Errorf("Argument `user` is mandatory")
	}

	redirects := make([]*redirectInfo, 0)
	for _, redirect := range flag.Args() {
		rinfo, err := newRedirectInfo(redirect)
		if err != nil {
			return err
		}
		redirects = append(redirects, rinfo)
	}

	for _, redirect := range strings.Split(os.Getenv("SSH_REDIRECTS"), "/") {
		if redirect == "" {
			continue
		}

		rinfo, err := newRedirectInfo(redirect)
		if err != nil {
			return err
		}
		redirects = append(redirects, rinfo)
	}

	if !*genBlueprint && len(redirects) == 0 {
		return fmt.Errorf("No redirects are defined, no reason to start.")
	}

	ssh, err := NewSshTunnel(*user, *host, *key, *blueprint)
	if err != nil {
		return err
	}

	defer ssh.Close()

	if *genBlueprint {
		fmt.Printf("Remote server key blueprint is: %s\n", ssh.KeyBlueprint())
		return nil
	}

	for _, r := range redirects {
		err = ssh.Redirect(r.Listen, r.Connect)
		if err != nil {
			return err
		}
	}

	return ssh.Wait()
}

func main() {
	err := mainErr()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
