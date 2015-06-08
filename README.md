Introduction
============

This is an example OAuth 2.0 utility that can:
 * Generate authorization URL, and retrieve access token and refresh token
 * Use a refresh token to exchange for an access token

This is not an official Google product.

Installation
============

To install:

    $ git clone https://github.com/saturnism/oauth2util.git
    $ cd oauth2util
    $ go get

To install it into `$GOBIN`

    $ go install

To run it as is:

    $ go run oauth2util.go

To build:

    $ go build oauth2util.go

To build for multiple platforms using Docker:

    $ script/build

Example Usage
=============

Client ID and Client Secret
---------------------------

You need these.  You must pass this in via the command line argument,
or alternatively, set the environmental variable.

Passing in via the command line:

    oauth2util --client_id=... --client_secret=... auth --scope=...
    oauth2util --client_id=... --client_secret=... exchange -t [refresh_token]

Using environmental variable:

    export OAUTH2_CLIENT_ID=...
    export OAUTH2_CLIENT_SECRET=...

Following examples assumes you've set the client ID and client secret into
environmental variables.

Authorization
-------------

Authorize with Google using the `profile` scope:

    oauth2util auth --scopes=profile

This will print out the token response, such as:

    {
      "access_token": "..."
      "token_type": "...",
      "refresh_token": "...",
      "expiry": "..."
    }

Authorize with different providers:

    oauth2util --auth_url=... --token_url=... auth --scopes=...

You can also set the authorization URL and token URL in environmental
variables:

    export OAUTH2_AUTH_URL=...
    export OAUTH2_TOKEN_URL=...
    oauth2util auth --scope=...

Refresh Token
-------------

You can exchange an existing refresh token for an access token:

    oauth2util exchange -t [refresh_token]

