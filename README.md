# XKCD_Project

This application retrieves the comics from the XKCD API, stores them in the database and then, displays the comics on UI.

## Table of Contents

- [Project Structure](#project-structure)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Contributing](#contributing)
- [License](#license)

## Project Structure

This project is a fullstack web application built using Golang. The frontend is rendered using Go Templates, and the backend is powered by Go with SQLite3 as the database. The application retrieves and displays comics from the XKCD API.

## Features

- Retrieve and display XKCD comics on the UI.
- Store fetched comics in an SQLite3 database.
- Simple, clean UI generated using Go Templates.
- Delete comics from the UI and from the database.

## Technologies Used

- **Backend:** Golang
- **Frontend:** Go Templates
- **Database:** SQLite3
- **External API:** XKCD API

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (version 1.22.1 or higher)

## Getting Started

Follow the steps below to set up the project locally.

### 1. Clone the Repository

```bash
git clone https://github.com/PatrickDiallo23/XKCD_Project.git
cd XKCD_Project
```

### 2. Install Dependencies

Install the necessary Go modules:

```bash
go mod tidy
```

### 3. Run the Application
Start the application using the Go command:

```bash
go run main.go
```

**Note**: This command will run the Terminal application where the comics are stored in memory and retrieved sequentially.
You can run this application can be run using flags.

- To run in sequential mode

```bash
go run main.go -mode sequential
```

- To run in concurrent mode

```bash
go run main.go -mode concurrent
```

- To run the fullstack application (xkcd mode)

```bash
go run main.go -mode xkcd
```

The application will be accessible at http://localhost:8080.

#### Usage

- Visit http://localhost:8080 to see the latest XKCD comic.
- Use the navigation buttons to browse through previous comics.
- The comics are stored locally in the SQLite3 database after they are fetched from the XKCD API.

## Contributing

Contributions are welcome! Please follow the standard GitHub workflow:

- Fork the repository.
- Create a new branch.
- Make your changes.
- Submit a pull request.

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
