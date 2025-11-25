Go-SSH-Tunnel
=============

Introduction
------------

This application allows to create SSH tunnel between two hosts. This is the same as use of SSH with autossh, but is especially useful for creating tunnels for your local deployments. It is the same as `cloudflared` but for self-hosted.

Usage
-----

This application can use only key pairs for authentication. Application supports 2 variants of defining the key pair:

1. Use `ssh-key-gen` to generate key pair and to copy id to the server to use this application.
2. Use `@secret` form of the key. In such case you can print corresponding public key using `gosshtunnel -key @key -showPublic` command.

Now you can easily run:

```bash
gosshtunnel -host ssh.example.net user username -key id_ed25519 -showBlueprint
```

This command will show public key blueprint for this server. This blueprint can be used (and I recommend to do it) to verify that server key was not changed. Now you can run:

```bash
gosshtunnel -host ssh.example.net user username -key id_ed25519 -blueprint blueprint-value 8000:8080 localhost:8888:localhost:2000
```

This command redirects remote port 8000 on all interfaces to localhost:8080 and localhost:8888 on server to localhost:2000 locally.

It is better to create custom SystemD service configuration if you need to create a tunnel from some public IP to your home.

For Docker you can use environment variables to set values. The same command using docker containter will be:

```bash
docker run -ti --rm -e USER=username -e SSH_REMOTE_HOST=ssh.example.net -e SSH_PRIVATE_KEY=id_ed25519 -e SSH_KEY_BLUEPRINT=blueprint-value -e SSH_REDIRECTS="8000:8080/localhost:8888:localhost:2000 gosshtunnel 
```
