# Project Daedalus Database Tool (PDT)

## Description

The Project Daedalus Database Tool (PDT) is a tool that allows users to interact with the Project Daedalus database. It provides a command-line interface that allows users to add, remove, and query data from the database. The tool is designed to be user-friendly and easy to use, with a simple and intuitive interface.

**NOTE: This tool requires specific authentication to be useful. Please contact the Author for more information.**

## Installation

The easiest way to install the Project Daedalus Database Tool is via Go:

```bash
go install github.com/DonovanMods/ProjectDaedalus-DB-Tool@latest
```

## Configuration

To use this tool, you need to configure it with your database connection details. You can do this by creating a configuration file in your home directory. The configuration file should be named `~/.pdtconfig.json` and contain the following data:

```json
{
  "firebase": {
    "credentials": {
      "YOUR FIREBASE JSON CREDENTIALS HERE"
    },
    "collections": {
      "meta": {
        "modinfo": "meta/modinfo",
        "toolinfo": "meta/toolinfo"
      },
      "mods": "mods",
      "tools": "tools"
    }
  },
  "github": {
    "token": "YOUR GITHUB ACCESS TOKEN HERE"
  }
}
```
