# Table Booking Platform Monorepo

## Development Environment

### Instructions

#### APIs

##### docker-compose

1. Navigate to an API directory and ensure certificates are present (see certificate generation).
1. Run `make dev` from an API directory to run an API using docker compose.
This will also run any database using docker-compose if needed.
   
---

##### local build

1. Navigate to an API directory and ensure certificates are present (see certificate generation).
1. Run `make local` from an API directory to build and run an API without dependencies.

##### certificate generation

APIs need a public and private key to communicate with each other. 
Run the following command to generate certificates. Place them in the individual API folders.

1. `openssl req -x509 -out localhost.crt -keyout localhost.key \
   -newkey rsa:2048 -nodes -sha256 \
   -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")`

#### User Interfaces

1. Go to the `booking_ui` lib directory.
1. Run `make dev` to install dependencies, generate graphql code and run user interface.

![](https://images.pexels.com/photos/1267708/pexels-photo-1267708.jpeg?cs=srgb&dl=four-women-sitting-on-benches-outside-building-1267708.jpg&fm=jpg)
