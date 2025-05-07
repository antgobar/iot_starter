# IOT Starter Project

An IoT platform starter project using Go, PostgreSQL, and soon TimescaleDb.

## Features

The user can create an register for an account. After logging in, the user can register their IOT device. The user will be given a temporarily visible API key which they will copy and use in their device. The device will send data to the server using the API key. The server will store the data in a PostgreSQL database. The user can view their devices and the data they have sent.

## TODO

- [X] Support monolith implementation using goroutines & channels
- [X] Allow future split into microservices with cmd apps
- [ ] Add a description of the project
- [ ] PostgresDb + TimescaleDb to store and retrieve data
