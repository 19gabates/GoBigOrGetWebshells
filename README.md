# Go Big or Get Webshells
GoBigOrGetWebshells (GBGW) is a specialized program designed to facilitate the creation and widespread deployment of webshells on various systems. Developed using the Go programming language, this program ensures seamless compatibility with both Linux and Windows operating systems. The webshells deployed through GBGW inherit the access privileges of the accounts on which they are installed.

It is crucial to note that GBGW is intended to be deployed by accounts with Admin/Sudo level access. Failure to deploy it with such accounts may result in the malfunctioning of multiple components of the program. Webshells, in this context, serve as a persistence method, enabling attackers to regain access to a system.

GBGW streamlines the deployment and concealment of webshells, simplifying the overall process and making it effortlessly efficient.

This project is the version 2.0 of the previous project [BiggerWebsites](https://github.com/19gabates/biggerWebsites) that only worked on Linux systems.

## Current Landmarks Reached
* Created Go Program to turn on a website with command execution
* Created Main Linux Program to handle Linux system deployment

## Future Landmarks
* Create a Main Windows Program to handle windows systems
* Work on techniques to hide a program on Linux and Windows
* Create a Linux kernal module to hide ports and PIDs

# Linux
In the Linux module of the Go program, the functionality is structured as follows:

1. Distribution Function:

      The initial function orchestrates the distribution of the website with command execution throughout the system, utilizing information specified in the config.txt file.

2. Binary Renaming Function:

      Subsequently, a second function comes into play, responsible for renaming the binary to a designated name as specified in the configuration. This strategic renaming aims to mitigate detection efforts, preventing adversaries from identifying the binary through a system-wide search.

3. Systemd Services Creation:

      Following this, a third function is implemented to establish systemd services. This ensures that, in the event of the process being terminated or the system restarting, the website automatically restarts, enhancing its resilience.

4. PID Concealment Function:

      The fourth function focuses on concealing the Process IDs (PIDs) of the websites. This deliberate obfuscation makes it more challenging for detection mechanisms to identify and track the running processes, thereby enhancing the overall stealth of the deployed websites.

## Example of Config File:
<Directory>,<Port>,<Filename>,<Service_Name>
```
/tmp,8080,BigWeb,website
/tmp/home,4444,WEB,Bigweb
```

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
