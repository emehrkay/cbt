# cbt

## Environment Variables

- RPC_PORT -- the port to run the rpc server, should be in the format of ":3333"
- PRIVATE_KEY_FILE -- the private key file
- PUBLIC_KEY_FILE -- the matching public key file

## Running Examples

Examples can be run via `make` or calling the executable directly

```
make ex1

// or
go run cmd/main.go ex1

// or build it and do the same as above
```

### Examples

*ex1*

logs in with a user, buys a ticket

*ex2*

logs in with a user, buys a ticket. Then logs in with admin and looks at the train details

*ex3*

logs in with a user and attemps to view train details, an admin-only action

*ex4*

logs in with a user and adds a ticket, logs in as admin and removes ticket from train

*ex5*

logs in with a user and adds a ticket, logs in as admin and updates the ticket. the user logs in and changes the ticket again. finally another user tries to unsuccessfully edit the first user's ticket

## Running Tests

```
make tests
```