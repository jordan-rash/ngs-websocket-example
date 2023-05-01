import Alpine from "alpinejs";
import { connect, consumerOpts, headers, JSONCodec, credsAuthenticator, Empty } from 'nats.ws';

Alpine.data("joke", () => ({
  nats: null,
  jc: null,
  joke: "",

  async init() {
    const ngs = "wss://connect.ngs.global"
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
`

    this.jc = JSONCodec()
    this.nats = await connect({servers: ngs, authenticator: credsAuthenticator(new TextEncoder().encode(creds))})
    
  },

  async newJoke() {
    await this.nats.request("new.joke", Empty, {timeout:3000})
      .then((m) => {
        this.joke = this.jc.decode(m.data).joke
      })
      .catch((err) => {
        console.log("error: "+ err);
      })
  },

}))

Alpine.start()
