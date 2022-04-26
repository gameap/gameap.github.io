---
title: GameAP Daemon specifications
layout: default
lang: en
category: GameAP Daemon
order: 60
---

GameAP Daemon components specification. Any component named as "server" in that document!

## Command server

Allows to run any commands (include os commands)

### Request

First list element [Binn](https://github.com/liteserver/binn) must contain command code
Second list element must contain command
Third list element must contain path to working directory

Example:

| Type      | Value                 | Description
|-----------|-----------------------|---------------------------
| uint8     | 1                     | OS command code
| string    | whoami                | Command
| string    | /srv/gameap           | Working directory


### Response

Response contains: gameap-daemon command code, os exit code, execution result

#### Response with data

| Type       | Description
|------------|---------------
| uint       | Code
| int32      | Exit status
| string     | Execution result

#### Коды ответов демона

| Code      | Description                                  
|-----------|-------------------------------
| 1         | Common error
| 2         | Critical error
| 3         | Unknown error
| 100       | Success

## File server

Server operate with data [Binn](https://github.com/liteserver/binn) List

### Request

First list element Binn must contain commend code (read dir, delete file, etc.)

Example (directory reading):

| Type      | Value                 | Description
|-----------|-----------------------|---------------------------
| uint8     | 4                     | Folder read code
| string    | /home/gameap          | Directory path
| string    | /home/gameap          | Directory path
| uint8     | 1                     | Reding mode (0 - files/directories list only, 1 - detail list with files info: size, modification date etc.)

### Response

#### Simple response

Simple response contains code (100 if success) and message.

| Type       | Description
|------------|---------------
| uint       | Code         
| string     | Message     

#### Response with data

Response with data (For example: directory listing)

| Type       | Description 
|------------|---------------
| uint       | Code           
| string     | Message     
| mixed      | Data        

#### Response codes

| Code       | Description                               
|------------|-------------------------------------------
| 1          | Common error
| 2          | Critical error
| 3          | Unknown error
| 100        | Success
| 101        | Ready to recieve/send file

### Files uploading/downloading 

#### Requiest

| Type      | Description                     | Value 
|-----------|---------------------------------|---------------
| uint8     | Code                            | 3
| uint8     | Mode (1 - download, 2 - upload) | 1
| string    | Path                            | /home/gameap/file.txt

#### Response

Server return simple response with command execution results
After answer "success", server will be ready to send/recieve files.

Server return answer on file recieve
| Type      | Description               | Value 
|-----------|---------------------------|---------------
| uint8     | Code                      | 101
| string    | Message                   | "File sending ready"
| uint64    | File size (bytes)         | 282345

### File listing

Get file listing.
There is 2 reading modes:
1) simle (only list)
2) detail (contains sizes/modification dates/etc.)

#### Request

| Type      | Description                   | Value
|-----------|-------------------------------|---------------
| uint8     | Directory reading code        | 4
| string    | Directory path                | /home/gameap
| uint8     | Mode (0 - simple, 1 - detail) | 1

#### Response

Data array

| Type      | Description
|-----------|-------------------------------------------
| string    | File name
| uint64    | File size (bytes)
| uint64    | Modification date (timestamp)
| uint8     | Type (1 - catalog, 2 - file, symlink, socket, pipe)
| uint16    | Permissions/chmod (example 755)

### Catalog creation

Create new catalog

#### Request

| Type      | Description               | Value
|-----------|---------------------------|---------------
| uint8     | Directory reading catalog | 5
| string    | Cirectory path            | /home/gameap/new_dir

#### Response

Server return simple response with command execution results.

### Move/Copy

Move or copy directory

#### Request

| Type      | Description                            | Value 
|-----------|----------------------------------------|---------------
| uint8     | Code                                   | 6
| string    | Current path                           | /home/gameap/old_file.txt
| string    | New path                               | /home/gameap/new_file.txt
| boolean   | true/false (false - move, true - copy) | true

#### Response

Server return simple response with command execution results.

### Remove

Remove file/directory

#### Request

| Type      | Description            | Value
|-----------|------------------------|---------------
| uint8     | Code                   | 7
| string    | Path                   | /home/gameap/delete_file.txt
| bool      | true/false (Recursive) | false

#### Response

Server return simple response with command execution results.

### Detail file information

Get detail information about file (cration/modification date, mime, etc.)

#### Request

| Type      | Description | Value 
|-----------|-------------|---------------
| uint8     | Code        | 8
| string    | Path        | /home/gameap/file.txt

#### Response

| Type      | Description                             
|-----------|-------------------------------------------
| string    | File ame
| uint64    | File size (bytes) 
| uint8     | Type <br> 1 - catalog;<br> 2 - file;<br> 3 - character device;<br>4 - block device;<br>5 - named pipe;<br>6 - symlink;<br>7 - socket;<br> 0 - unknown
| uint64    | Modification date (timestamp)
| uint64    | Access date (timestamp)
| uint64    | Creation date (timestamp)
| uint16    | Permissions (chmod)
| string    | Mime

### Change permissions

Change file permissions

#### Request

| Type      | Description         | Value 
|-----------|---------------------|---------------
| uint8     | Code                | 9
| string    | Path                | /home/gameap/file.txt
| uint16    | Permissions (chmod) | 0755

#### Response

Server return simple response with command execution results.

## Status server

Allows to recieve some info about daemon (tasks list, servers online, etc.)

### Request

First list element Binn must contain command code

Example:

| Type      | Value         | Description             
|-----------|---------------|---------------------------
| uint8     | 2             | Base daemon info


### Response

#### Response with data

Response with data (example directory listing)

| Type      | Description 
|-----------|---------------
| uint      | Code           
| string    | Message     
| ...       | Data        
| ...       | Data        
| ...       | Data        

#### Response codes

| Type      | Description                               
|-----------|-------------------------------------------
| 1         | Common error
| 2         | Critical error
| 3         | Unknown error
| 100       | Success

### Daemon version

Recieving info about gameap version and compilation date.

#### Request

| Type      | Description              | Value 
|-----------|--------------------------|---------------
| uint8     | Code                     | 1

#### Response

| Type      | Description 
|-----------|---------------
| uint      | Code
| string    | Message
| string    | GameAP Daemon version
| string    | Compilation date and time

### Base working data

#### Request

| Type      | Description              | Value
|-----------|--------------------------|---------------
| uint8     | Code                     | 2

#### Response

| Type      | Description
|-----------|---------------
| uint      | Code
| string    | Message
| uint32    | Uptime GameAP Daemon
| uint32    | Running tasks count
| uint32    | Waiting tasks count
| uint32    | Online servers count

### Detail working data

#### Request

| Type      | Description              | Value
|-----------|--------------------------|---------------
| uint8     | Code                     | 3

#### Response

| Type      | Description 
|-----------|---------------
| uint      | Code
| string    | Message
| uint32    | GameAP Daemon uptime
| list      | Running tasks ID list
| list      | Waiting tasks ID list
| list      | Online servers ID list
