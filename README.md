# catalyze-go
Simple REST API written in Go

## Setup

The following Go dependencies need to be installed:

    go get -u github.com/gorilla/mux
    go get -u github.com/urfave/negroni
    go get -u github.com/go-sql-driver/mysql
    go get -u github.com/satori/go.uuid
    go get -u golang.org/x/crypto/bcrypt

Once you have those installed, run the `db.sql` file on a MySQL database named "catalyze" and setup a username "catalyze" with a password of "abcd1234". I don't have a neat way to change the settings yet, unfortunately.

## Usage
Make sure to include the `Authorize: <token>` header for all requests except for `POST /user`, `POST /auth`, and `GET /`.

### GET /
- Accessing without the Authorize header will give you a "Hello World" message
- Accessing with a valid Authorize token will give you the user data

### POST /auth
Create a new access token. Tokens are valid for 24 hours by default.

Example Input: 

    {
      "username": "brubbles",
      "password": "barneyrules"
    }

Example Output: 

    {
      "token": <token>
    }  
    
### DELETE /auth
Invalidates access token that is in the Authorization header. No input necessary.

### POST /user
Creates a new user.

Example Input: 

    {
      "username": "yabadabadood",
      "password": "wilma4",
      "firstname": "Fred",
      "lastname": "Flintstone",
      "color": "orange"
    }

### GET /user/{username}
Access the user data of the given username. Authorized users can only only look at themselves, so this would be useful if there was support some sort of admin user. 

Example Output: 

    {
      "id": 2,
      "firstname": "Wilma",
      "lastname": "Flintstone",
      "username": "cavelady",
      "color": "white"
    }
    
### PUT /user/{username}
Updates the user with the given username. Authorized users can only update themselves for now.

Example Input: 

    {
      "firstname": "Betty",
      "lastname": "Rubbles",
      "color": "blue"
    }
    
### DELETE /user/{username}
Deletes the user with the given username. All associated access tokens automatically get deleted as well. Authorized users can only delete themselves for now.
