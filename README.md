## Steps

### Create an NGS Account
Head over to the [NGS site](https://app.ngs.global/) and create a free account.
Follow the instructions to get `nsc` configured locally with your free account credentials

### Generate the creds for wasmCloud
First, create a user for your new NGS account 
`nsc create user -n websocket`

Second, generate the credential file that you will need for your client 
`nsc generate creds -n websocket -o ./websocket.creds`

### Use the `nats.ws` project to connect to NGS

Clone this repo down and edit the `ui/main.js` file where you see the placeholder for NGS creds (see below).
You will copy and paste the output of the `nsc generate creds` command above
```javascript
const creds = `-----BEGIN NATS USER JWT-----
<JWT>
------END NATS USER JWT------

************************* IMPORTANT *************************
NKEY Seed printed below can be used to sign and prove identity.
NKEYs are sensitive and should be treated as secrets.

-----BEGIN USER NKEY SEED-----
<NKEY>
------END USER NKEY SEED------

*************************************************************
```
### Build the two actors
This repo includes 2 actors.  One holds the "business logic" for getting a joke and the other holds the UI.

For the UI directory:
```
cd ui/
go generate
tinygo build --target wasi .
wash claims sign --name ui ws.wasm --http_server -l
```

For the joke directory:
```
cd joke/
tinygo build --target wasi .
wash claims sign --name ui ws.wasm --http_client -m
```
### Start all the pieces
You will need to start 3 providers from the wasmCloud Azure Registry

- httpserver
- httpclient 
- nats_messaging

And the two actors we just build: 
- ui 
- joke

There will be three links:
- ui <-> httpserver 
- joke <-> httpclient 
- joke <-> nats_messaging 

The link variables for the nats_messaging should look similiar to this 
`CLIENT_JWT=<JWT FROM nsc CREDS>,CLIENT_SEED=<SEED FROM nsc CREDS>,SUBSCRIPTION=new.joke,URI=connect.ngs.global `
