# SU-Organization-Matching

Repo for the SU Organization Matching Capstone Project

## Initial project proposal

Credit to Dr. Debika Sihi and her students for the initial proposal

"Multiple students have shared that they didn’t learn about many student organizations until their senior year. They suggested that an automated system matching students to organizations based on their interests (pulled from their application or survey responses) would be incredibly helpful. As first-year students, many expressed feeling overwhelmed and unsure of where to even begin looking for community."

## Team

This team is is comprised of Tanner Klein (TannerK7), Matthew Volkin (MattVolkin), Ben McKallip (bmck039), and Aidan Balakrishnan (drumb0y).

### Contributions
- Tanner Klien: Project manager, sorting logic and quiz question design
- Matthew Volkin: UI/UX and API integration 
- Ben McKallip: Backend and Database
- Aidan Balakrishnan: UI/UX

## SSH Setup

### Setup Script
In a command prompt run:

Windows: `ssh_setup.bat`

Linux: `bash ssh_setup.sh`

These scripts will install and configure the necessary connector software. It will open a web browser and prompt for a login, simply enter your SU email, then the code it emails you, then grant access. 

### Connecting 
* To use the updated server, the ssh key needs to be reset. To accomplish this, run the command: `ssh-keys -R capstonessh.benmckallip.com`
  * If this doesn't work, simply remove the `known_hosts` file from the `$USERNAME$/.ssh` folder
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

## Running
### Build Frontend
In the "Svelte Examples/plain-svelte-app" directory, run the following:
```
npm run build
```
### Run Server
In the "Server Examples" directory, run this:
```
go run server.go
```

## Hosting:
To view our product go to https://capstone.benmckallip.com/

## Files: 

- API_REFERENCE.md: Examples, usage guide, and syntax information about each API method, what parameters it takes, and what given variable/category names are.
- build_and_run.bat: Executable file to build server.
- build_and_run.sh: Shell runnable file to build server.
- dbout.txt: Text representation of true false values to rest swiping application.
- go.mod: List of direct libraries or tools used by the go portions of the application (especially for API calling).
- go.sum: List of all libraries and their dependencies for the go portions of the application.
- init_db.sql: Initialization for the database structure.
- package-lock.json: Node.js library dependencies for application.
- package.json: Node information for this application.
- README.md: README that illustrates project goals, process, and other information.
- sorting.go: Go file that is used to test out different (hypothetical) users for testing the matching algorithm to refine it.
- ssh_setup.bat: SSH setup file for developers to be able to ssh into a communal machine to work on code together and have the application webpage update in real time.
- ssh_setup.sh: SSH setup file for developers to be able to ssh into a communal machine to work on code together and have the application webpage update in real time.
- user_setup.sh: Setup the various development accounts with the correct permissions to edit and execute files.
- vite.config.js: Setup the settings for the webpages used for this project and establish the routings.

### BaseCSVTemplate
- Question Planning: List of Adjectives.csv: List of adjectives used primarily by the swiping and matching parts in a CSV format to show the formatting that is used for the database.
- Question Planning: Reformatted Data.csv: Formatted CSV that reflects the different organizations and what preferences/traits they are aiming towards.
- Results Page.csv: Names, times, descriptions, etc for each organization in the system in a csv format.

### cmd
#### import_organizations
- main.go: Read and parse through Question Planning: Reformatted Data.csv file and add organizations to DB.

#### import_questions
- main.go Read and parse through Question Planning: List of Adjectives.csv file and add question (adjectives) to DB.

### Components
- AdminCreateNewClub.svelte: Component that allows a user with an admin account (one of the projects creators or a faculty member of Student activities) to create a new club on the application that can be populated with information later by the respective clubs officers.
- AdminHome.svelte: Accessible only to the admins, this component provides the above page the power to edit any already existing organizations.
- AdminSwitch.svelte: Component that allows for an Admin to switch into officer or base user views for testing and checking purposes.
- APIHandler.svelte: Component that does not do any visual displaying but makes sure that any API calls are all routed the same way, through the element to streamline calls and dataflows.
- Card.svelte: Cards for each adjective/activity in the swiping portion of the application, these read directly from the database to generate the correct amount of cards as well as randomize their order.
- deleteAccount.svelte: Component that performs checks that the user wants to consciously take this action and sends API calls to clear an accounts data.
- Footer.svelte: Element that displays the team that made this project as well as purpose. At the bottom of every webpage.
- Header.svelte: Element that displays the name of the project, sign in/out option, and available pages for a user’s navigation.
- Login_popup.svelte: Component that prompts the user to sign into the google account service (and returns a secure token so that our application can accurately recall their information as needed) to match their information to an organization.
- Mult_choice_demo.svelte: Demo to show how the multiple choice components work.
- ProgressBar.svelte: Custom component that can be used for wherever is needed (including the swiping component) to indicate to the user how close they are to the end.
- Results.svelte: Display a club’s name, image of the club, associated activities, and much more information to the user. Also takes care of necessary API calls.
- SettingPage.svelte: Handles API calls and other information to allow for an officer to manage/edit their organization.
- SwipingApp.svelte: Part of the quiz, allows the user to swipe left or right depending on if they like a topic or activity.

#### OrgPhotos
- (photos for each organization).jpg: Each organization provided a photo to represent them on the final homepage.

#### pages
- AdminCreateNewClub.svelte: Page that allows a user with an admin account (one of the projects creators or a faculty member of Student activities) to create a new club on the application that can be populated with information later by the respective clubs officers.
- AdminHomeFinal.svelte: Accessible only to the admins, this page allows for the overview and editing of any already existing organizations.
- DeleteAccountPage.svelte: Page that allows for a basic user (like a student) to clear all of their preferences/reset their account of they want their information cleared from our system.
- DemographicQuiz.svelte: Page that handles the user interface that allows them to put in their demographic information (gender, ethnicity, religion, etc) as those can be factors that help recommend aligned orgs.
- HomePage.svelte: Homepage that proposes the concept of the project as well as the option to start the demographics quiz. First page a user typically encounters upon encountering our site.
- Manual.svelte: Displays the user manual.
- ResultsPage.svelte: After taking both the demographics quiz as well as the swiping portion of the project, this component displays the best organizations for each user sorted by their closest match.
- SettingsPage.svelte: For officers or above, this page allows for the user to change aspects of a given organization.
- SwipingPart.svelte: This page allows for the user to show what aspects of clubs they could be interested in in a swipe left/right format.

### matching
- sorter.go: Go file that contains the algorithms that evaluates a user and sorts them according to their likes, dislikes, and demographic information into a descending list of their most aligned organizations.

### node_modules
- (...) various libraries we use

### server

- db.go: Access the DB (using ent) as an object to make requests.
- login.go: Allow for the user to use the google profile authentication and the DB to store their unique tokenID for later recalling of data/results.
- server: Compiled binary file of server.go.
- server.go: Top level server file that builds and compiles the server and allows for the connectivity between the DB, front-end webpages, with correct permissions given a user’s type of access.
- server.log: Created when the build script is run to log any activities that are happening within the project.

#### ent
- (...).go: Various ent generated files that allow for the user to treat the database as an Object with methods and access it in a more user friendly way, with ent formatting any SQL queries. 


###### schema
- answer.go: Documents answer-table database schema for ent
- club.go: Documents club-table database schema for ent
- question.go: Documents question-table database schema for ent
- user.go: Documents user-table database schema for ent

### Webpages

- (...).html: User-generated html files that have a basic representation of each front-facing webpage

#### dist (Files in this folder were automatically created by the build methods, not user generated)

- (...).html: Html parts created in the build process for each webpage 

##### assets
- (...).js: Javascript part for each front-facing webpage
- (...).css: CSS parts for each front-facing webpage

#### src
- mountPage.js: mountPage method to call a svelte function from JS

##### entries
- (...).js: Using the mount function in svelte, these files (one for each webpage) mount the files from the component folder
