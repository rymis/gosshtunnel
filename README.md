Go-SSH-Tunnel
=============

Introduction
------------

This application allows to create SSH tunnel between two hosts. This is the same as use of SSH with autossh, but is especially useful for creating tunnels for your local deployments. It is the same as `cloudflared` but for self-hosted.

Usage
-----

You need to generate key-pair using `ssh-key-gen` and to copy id to the server to use this application. Then you can easily run:

```bash
gosshtunnel -host ssh.example.net user username -key id_ed25519 -showBlueprint
```

This command will show public key blueprint for this server. Then

```bash
gosshtunnel -host ssh.example.net user username -key id_ed25519 -blueprint <blueprint-value> 8000:8080 localhost:8888:localhost:2000
```

This command redirects remote port 8000 on all interfaces to localhost:8080 and localhost:8888 on server to localhost:2000 locally.

It is better to create custom SystemD service configuration if you need to create a tunnel from some public IP to your home.

For Docker you can use environment variables to set values.
