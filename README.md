<img align="right" width="33%" style="margin-bottom: 3em" src="./logo.svg">

# Hi, we're Labonion! :wave:

Welcome to UMS by Labonion! This repository is meant to be a starting point for Organizations who want to build their backend using go lang and mongodb.

### And by “organizations”, we mean a _company_ with alot of people or a _startup_ with very few...

### We :blue_heart: open source, y'all.

This project evolved from one of my personal project that I am working on which cannot be open sourced yet.

### Pre Requisites

1. Install Go version >1.21.6 

### Run Server

1. Copy .env from .env.example
```make copy_env```

2. Run the server
```make run_server```

### What does this repository do?
1. This is a User Management System built using golang.
2. We use Redis for JWT Authentication.
3. This backend uses mongodb as its primary database.
4. All services are containerized using docker.
5. It uses a common repository interface that allows you to create controllers on the go.
6. It has pre built redis and mongodb connections for seamless development
7. For Authentication all you need is a header with X-API-KEY and the UUID Generated after login.
8. The AI Module exposes an api where it can stream output based on prompt
9. Spaces - Spaces is a concept where you can add mutiple users to a space and converse, manage and monitor within that space.

### But most of all?

We're excited to see what you build with UMS- by Labonion! :owl:

We are also building something of our own,

:onion:   [Markie - by Labonion][labonion] on ReadMe. <br>

[labonion]: https://labonion.com