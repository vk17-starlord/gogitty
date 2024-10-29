

# GoGitty: A Git-Like Version Control System in Go

Article - [Write yourself a git](https://wyag.thb.lt/)

GoGitty is a simplified, Git-inspired version control system built from scratch in Go. This project serves as a learning tool to understand Git internals by re-implementing essential Git commands. By using the Cobra CLI framework, GoGitty provides a command-line interface that mimics the basic functionality of Git, such as initializing repositories, adding and committing changes, and viewing logs.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Folder Structure](#folder-structure)
- [Commands](#commands)
- [Contributing](#contributing)

## Features

GoGitty currently will supports the following Git-inspired commands:

- **init**: Initialize a new repository
- **add**: Stage files to be committed
- **commit**: Commit staged changes to the repository
- **log**: View commit history
- **checkout**: Switch branches or restore working tree files
- **tag**: Tag commits with meaningful names
- **status**: Show the working directory and staging area status
- ...and more!

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or higher
- [Cobra CLI](https://github.com/spf13/cobra) for generating CLI commands

### Clone the Repository

```sh
git clone https://github.com/yourusername/gogitty.git
cd gogitty
```

### Install Dependencies

Make sure Cobra CLI is installed for command generation, if you plan on adding new commands:

```sh
go install github.com/spf13/cobra-cli@latest
```

### Build the Project

Compile the project into an executable:

```sh
go build -o gogitty main.go
```

This creates a `gogitty` executable in your project directory.

## Usage

Initialize a GoGitty repository in any directory by running:

```sh
./gogitty init
```

From here, you can stage files, commit changes, view the log, and explore other commands.

```sh
./gogitty add <file>
./gogitty commit -m "Commit message"
./gogitty log
```

## Folder Structure

Here's a breakdown of the project structure for easy navigation and contribution:

```plaintext
gogitty/
├── cmd/                # CLI command definitions
│   ├── root.go         # Root command
│   ├── add.go          # 'add' command implementation
│   ├── init.go         # 'init' command implementation
│   └── ...             # Other command files
├── internal/           # Internal packages for core functionality
│   ├── core/           # Core Git logic (repo, objects, etc.)
│   ├── storage/        # Storage and encoding/decoding logic
│   ├── refs/           # References (branches, tags) management
│   └── common/         # Shared utilities and constants
├── pkg/                # Exported packages, like custom logger (optional)
├── scripts/            # Setup and test scripts
├── main.go             # Entry point of the application
└── README.md           # Project documentation
```

## Commands

### Init

```sh
./gogitty init
```

Initializes a new GoGitty repository in the current directory.

### Add

```sh
./gogitty add <file>
```

Stages a file for the next commit. The file is added to the staging area.

### Commit

```sh
./gogitty commit -m "Commit message"
```

Commits staged files to the repository with a commit message.

### Log

```sh
./gogitty log
```

Displays the commit history of the repository.

### Other Commands

- **cat-file**: Shows the content of an object by its hash.
- **status**: Shows the status of files in the repository.
- **tag**: Tags a specific commit with a label.

Refer to individual command help for more options:

```sh
./gogitty <command> --help
```
