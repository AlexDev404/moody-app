# moody

This project seeks to satisfy the requirements of Test #1 in the [CMPS3162 - Advanced Databases](https://doit.ub.edu.bz/course/view.php?id=70) course

## What this project is

This is a webapp/microfrontend that analyzes a user's mood (based on text input, emoji selection, or even a selfie analysis) and generates a personalized music playlist using AI.

## Prerequisites

You are required to have the following dependencies installed on your system before running this application:

- [PostgreSQL](https://www.postgresql.org/download/)
- [Go](https://go.dev/)
- [Golang-Migrate](https://github.com/golang-migrate/migrate)
- [Node 22.x](https://nodejs.org)

## To run

> [!IMPORTANT]
> On First Run (If not, skip to below)
> Open a command-line on either Windows or Linux inside of the project root and execute the following
>
> ```shell
> # Install all dependencies and scaffolding required
> $ make prepare
> $ make initialize
> # Run the migrations
> $ make migrate

Please ensure you also create and edit the `.envrc` before running any of these commands.
There is an example file included in the project you can use to base it off of.

To run the project, you can execute the following each time:

```shell
$ make
"Now listening on port http://127.0.0.1:4000"
```

## What it looks like

|Explanation|Preview|
|:-----------|:-------:|
|Homepage|![image](https://raw.githubusercontent.com/AlexDev404/moody-app/refs/heads/main/docs/Screenshot_2025-04-13_201419.png)|
|Tools page|![image](https://raw.githubusercontent.com/AlexDev404/moody-app/refs/heads/main/docs/Screenshot_2025-04-13_201444.png)|

## Database Diagram

![database diagram - image](https://raw.githubusercontent.com/AlexDev404/moody-app/refs/heads/main/docs/database_diagram.png)

## Technologies used

- This project makes use of Go's `net/http` while statically serving the stylesheets.
- It also makes use of Web Assembly for some of the interactions
- This project also makes use of [Web Components](https://developer.mozilla.org/en-US/docs/Web/API/Web_components) and the [Shadow DOM](https://developer.mozilla.org/en-US/docs/Web/API/Web_components/Using_shadow_DOM) for separating views and components effectively
- This project also makes use of [`popstate`](https://developer.mozilla.org/en-US/docs/Web/API/Window/popstate_event) for routing in a SPA-like manner (no reloading) and for the application lifecycle
- This project also makes use of `html/template` for the rendering
