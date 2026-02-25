# SU-Organization-Matching

Repo for the SU Organization Matching Capstone Project

## Initial project proposal

Credit to Dr. Debika Sihi and her students for the initial proposal

Multiple students have shared that they didn’t learn about many student organizations until their senior year. They suggested that an automated system matching students to organizations based on their interests (pulled from their application or survey responses) would be incredibly helpful. As first-year students, many expressed feeling overwhelmed and unsure of where to even begin looking for community.

## Team

This team is is comprised of Tanner Klein (TannerK7), Matthew Volkin (MattVolkin), Ben McKallip (bmck039), and Aidan Balakrishnan (drumb0y)

While all members will be working on various portions of the project together, the general roles are as follows

Leadership/Project Management - Tanner
Database Management - Ben
UX/UI - Matt, Aidan
Communication with other student orgs - Aidan and Tanner

## Task list

- [ ] Create a webpage (using React?) to start showing concepts of how we want the user to interact with entering their information (Aidan and Matt)
- [] Start building and researching the best way to set up a database, as well as what kind (Ben)
- [x] Finish form for club officers to fill out to give data about the individual organizations across campus (Tanner)
- [ ] Create brand guidelines document that discusses what colors and logos to use, as to fit into the Southwestern image (all)

## SSH Setup

### Setup Script
In a command prompt run:

Windows: `ssh_setup.bat`

Linux: `bash ssh_setup.sh`

These scripts will install and configure the necessary connector software. It will open a web browser and prompt for a login, simply enter your SU email, then the code it emails you, then grant access. 

### Connecting 
* Username is your name (lowercase). so for Ben, the username is `ben`. You can connect with `ssh ben@capstone` (replacing `ben` with your username)
* first-time password should be `ubuntu`, it will prompt you to change this
* `~/capstone` is a shared folder with much more storage than the home directory, try to keep shared files in there. 
* I have the server set up to host localhost:8080 at [capstone.benmckallip.com](capstone.benmckallip.com), so we can run the server and visit that site to test any changes.


## Installing
### Prerequisites

The project requires the following software versions:

Go: 1.23+

PostgreSQL: 16+

Node.js: 24+ (via NVM)

### Linux (Debian-based)

1. Install System Packages

Install the core dependencies using apt:
```
sudo apt update
sudo apt install golang-1.23-go golang-1.23-src postgresql
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.4/install.sh | bash
```

2. Configure Node.js

Use NVM to install the specific Node version:
```
nvm install 24
nvm use 24
```

3. Configure PostgreSQL

Edit the configuration file to allow local connections:
```
sudo nano /etc/postgresql/16/main/postgresql.conf
```

Find and uncomment: listen_addresses = 'localhost'

4. Database Setup

Access the PostgreSQL prompt to create the user and database:
```
sudo -u postgres psql
```

Run the following SQL commands:
```SQL
CREATE USER dev_user WITH PASSWORD 'testing';
CREATE DATABASE dev_project_db OWNER dev_user;
GRANT ALL PRIVILEGES ON DATABASE dev_project_db TO dev_user;
\q
```

Restart the service to apply changes:
```
sudo systemctl restart postgresql
```

### Windows

1. Install System Packages

Open PowerShell as Administrator and run:

#### Install Go and PostgreSQL
```
winget install GoLang.Go --version 1.23
winget install PostgreSQL.PostgreSQL.16
```

#### Install NVM for Windows (requires a new terminal after install)
```
winget install CoreyButler.NVMforWindows
```

2. Configure Node.js

Open a new PowerShell window and run:
```
nvm install 24.0.0
nvm use 24.0.0
```

3. Database Setup

PostgreSQL on Windows usually listens on localhost by default. First, set a password for the postgres superuser if you haven't already:

Open Command Prompt (not PowerShell) as Administrator and run:
```
"C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -c "ALTER USER postgres WITH PASSWORD 'postgres_password';"
```

If using PowerShell instead, run:
```powershell
& "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -c "ALTER USER postgres WITH PASSWORD 'postgres_password';"
```

Then, open the SQL Shell (psql) from the Start Menu. When prompted:
- **Server [localhost]**: Press Enter (default is fine)
- **Database [postgres]**: Press Enter (default is fine)
- **Port [5432]**: Press Enter (default is fine)
- **Username [postgres]**: Press Enter (default is fine)
- **Password**: Enter the password you just set (or the one set during installation)

Then run the following SQL commands:
```SQL
CREATE USER dev_user WITH PASSWORD 'testing';
CREATE DATABASE dev_project_db OWNER dev_user;
GRANT ALL PRIVILEGES ON DATABASE dev_project_db TO dev_user;
\q
```

## Project Initialization

After setting up the environment, run these commands in the project root folder:

Initialize Database Schema

#### Provide password 'testing' when prompted
```
psql -h localhost -U dev_user -d dev_project_db -f init_db.sql
```

Backend Setup (Go)
```
cd 'Server Examples'
go mod tidy
cd ..
```

Frontend Setup (Svelte)
```
cd 'Svelte Examples/plain-svelte-app'
npm install
```

## Running:
### Build Frontend:
In the "Svelte Examples/plain-svelte-app" directory, run the following:
```
npm run build
```
### Run Server:
In the "Server Examples" directory, run this:
```
go run server.go
```