# Go Big or Get Webshells
GoBigOrGetWebshells (GBGW) is a specialized program designed to facilitate the creation and widespread deployment of webshells on various systems. Developed using the Go programming language, this program ensures seamless compatibility with both Linux and Windows operating systems. The webshells deployed through GBGW inherit the access privileges of the accounts on which they are installed.

It is crucial to note that GBGW is intended to be deployed by accounts with Admin/Sudo level access. Failure to deploy it with such accounts may result in the malfunctioning of multiple components of the program. Webshells, in this context, serve as a persistence method, enabling attackers to regain access to a system.

GBGW streamlines the deployment and concealment of webshells, simplifying the overall process and making it effortlessly efficient.

This project is the version 2.0 of the previous project [BiggerWebsites](https://github.com/19gabates/biggerWebsites) that only worked on Linux systems.

## Current Landmarks Reached
* Created Go Program to turn on a website with command execution
* Created Main Linux Program to handle Linux system deployment


## Linux Webshell Usage

```
Usage:
      -all      : Run all functions in order
      -ldist    : Run distribution function
      -lname    : Run renaming function
      -lcreates : Run create systemd service and reload systemd
      -lstart   : Reload systemd and enable/start services
      -lhide    : Ensure empty PID file, update PID file, and hides PIDs
```
