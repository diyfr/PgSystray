# PostgreSQL Portable Launcher for Windows 
[Initial project](https://github.com/bvp/PostgreSQLPortable)   

## Description
It allows you to control the portable version PostgreSQL from the tray icon.  

## Building for Windows  
You need `windres` for resource generation. Available from [MinGW](https://sourceforge.net/projects/mingw-w64/files/) installation.  
Add `[YOUR FOLDER]\mingw64\bin`  in your PATH variable   
Missing library :  
```
go get github.com/PuerkitoBio/goquery
go get github.com/andybalholm/cascadia
go get github.com/blang/semver
go get github.com/lxn/walk
go get github.com/lxn/walk/reflective
```

Build  
```
go generate
go build -ldflags "-H=windowsgui"
```

## Current status
 - [x] Windows support
 - [ ] Linux support
 - [ ] OSX support
 - [x] Config file
 - [x] Download PostgreSQL distributive
 - [x] Extracting downloaded archive
 - [x] Custom Port 
 - [ ] Show settings dialog on first start
 - [ ] Download/Launch pgAdmin
 

