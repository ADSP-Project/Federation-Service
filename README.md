# Federation Service Setup

This repository contains the code for setting up a federated marketplace using Go and PostgreSQL.

## Prerequisites

- Go

## Getting Started

1. Clone the repository:

   `git clone https://github.com/ADSP-Project/Federation-Service.git`

   `cd Federation-Service`

2. Install dependencies:

   - Make sure to initialize go.mod to manage dependencies:
   
      `go mod init github.com/ADSP-Project/Federation-Service`

   - Fetch and arrange them into newly generated go.mod:
   
      `go mod tidy`

   - Finally, install the required dependencies with the following command:
   
      `go mod download`

3. Simulate a shop joining the federation:
   - To simulate a shop joining the federation, open a new terminal and run the following command:

     `go run shop.go [port] [name]`

     Replace `[port]` with the desired port number and `[name]` with the name of the shop.

     This will start a shop server that automatically joins the federation by sending a POST request to the federation server.

     **Important:** Hub and AuthServer from Federation-Hub should be running so that shop can join Federation.

4. Additional Notes:
   - You can run multiple instances of the shop server by providing different port numbers and shop names.
