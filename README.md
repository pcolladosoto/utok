# uTok: A microscopic OpenID connect client
uTok aims to be a micro ($\mu$) client for generating access tokens through OpenID Connect's
Device Authorization Flow.

As hinted by its default settings, `uTok`'s main target is the
[WLCG's main Indigo IAM instance](https://wlcg.cloud.cnaf.infn.it). You can find a bit
of documentation on its APIs and such [here](https://indigo-iam.github.io/v/current/).
Even though we haven't tested it, `utok` might work with other *issuers*: we didn't really
do anything 'special' for targetting Indigo IAM when it comes to token generation.

Bear in mind the official client for Indigo IAM is [`oidc-agent`](https://github.com/indigo-dc/oidc-agent),
but we found it a bit 'aggressive' in its pursue of `ssh-agent`'s behavior and, after a lot of digging,
we didn't manage to get it to work on newer macOS versions or on CentOS 7...

This client's interface is rather self explanatory: running `utok` with no arguments will show some
pointers to make use of `utok`.

## Installation
You can just download the latest build for your platform and place the binary anywhere on your `PATH`.

Uninstalling `utok` is a matter of removing that binary!

## Getting the first token
In order to get a token you first need to create a client:

    $ utok cli create

This will create `~/.utok/client.json` containing the reply's contents. This reply will
also be shown on screen.

After that, you can generate tokens with:

    $ utok token

This instructs `utok` to read the contents of `~/.utok/client.json` to then try to generate a
token. If none have been generated previously, the Device Authorization Flow will be triggered
so that you'll need to navigate to a particular URL and input a code: all these instructions
will be shown on screen. The generated token will be stored on `~/.utok/token.json`.

After generating the first token, `utok` will leverage the **refresh token** embedded in the
initial one to re-generate access tokens at will. However, this is completely transparent:
the user need only run `utok token`. Fresh tokens will be stored on `~/.utok/token_fresh.json`.

After a client is no longer needed, it can be deleted with:

    $ utok cli delete

And... that's it really! Happy tokening!
