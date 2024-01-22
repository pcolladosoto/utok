% utok(1) | General Commands Manual

# NAME
utok - Micro OpenID Connect client

# SYNOPSIS
`utok [-h | --help] [cli | config | help | token | version]`

# DESCRIPTION
The utok client provides an easy and convenient manner of generating access tokens through the
OpenID Device Authorization Flow. It's been designed specifically targeting the Indigo IAM
service, but it should work with standards-compliant IAM implementations.

Aside from token generation, utok will also generate/register OpenID clients through which it
then generates tokens. The capability of deleting clients is also offered.

The different JSON documents represnting both the clients and tokens will be stored under
`~/.utok` so that they can be queried at will. Given the sensitive information contained in
these documents, permissions are set up so that only the owner can query the contents of
these different files.

Bear in mind utok's help message (accessible with `utok [--help | -h | help]`) contains a wealth
of information that's worth checking out if in doubt.

Happy tokening!

# OPTIONS
`-h, --help`

:   Show the help message and exit.

`--cli-contacts`

:   Comma-separated OIDC client contacts. Default: "foo@faa.com".

`--cli-name`

:   Name of the OIDC Client to create. Default "uTok-cli".

`--iss`

:   The issuer URL to contact. Default: "https://wlcg.cloud.cnaf.infn.it/".


`--scopes`

:   Comma separated scopes to request when generating tokens. Default "storage.read:/atlasdatadisk/SAM/,openid,offline_access".

# COMMANDS
utok's design leverages a series of commands, each taking care of a specific and distinct objective. One can query information
on a command by running `utok <command> --help`.



# CONFIGURATION
Just like the report, configuration is defined through a JSON file. This configuration is then made
available to the program through the `--conf` option. The available fields are:

**description [string]**

:   The description to be embedded (as-is) in the output report.

**interval [number]**

:   The delay, in seconds, between subsequent SNMP queries to the border switches.

**borderSwitches [array of borderSwitch]**

:   An array of the different `borderSwitch`es to monitor. These will be queried over SNMP
    every `interval` seconds.

**borderSwitch [object]**

:   A `borderSwitch` defines both the interfaces to monitor over SNMP and the SNMP parameters
    (i.e. *hostname*, *version* and *community*) necessary for gathering the necessary data.

**borderSwitch.hostname [string]**

:   The hostname of the border router. This will be resolved through an **A DNS** query to an
    **IPv4** address with which to interact over **UDP** on port **161**. A raw IPv4 address
    can also be provided here, but then the output report will identify this switch's interfaces
    with that IPv4 address instead of a qualified hostname.

**borderSwitch.publishedHostname [string]**

:   The hostname to publish associated to this border switch. Some sites might enforce security
    policies preventing them from publishing real (i.e. resolvable) hostnames for switches to
    avoid the dissemination of information regarding the real networking infrastructure. If left
    empty, the actual switch hostname will be published. This string does not have to be related
    in any way whatsoever with the border switch's real hostname.

**borderSwitch.hcSupport [bool]**

:   Whether the switch supports *High Capacity* (HC) counters. This affects the OID being
    requested over SNMP. If HC support is present, the `ifXTable` (OID 1.3.6.1.2.1.31.1.1)
    will be queried, otherwise, the `ifTable` (OID 1.3.6.1.2.1.2.2) will be queried instead.

**borderSwitch.snmpVersion [string]**

:   The SNMP version to use. This option **must be** one of `"1"`, `"2c"` or `"3"`. However,
    bear in mind there is no support for SNMPv3-specific options at the moment.

**borderSwitch.snmpCommunity [string]**

:   The SNMP community to use. This will usually be provided by the site's network administrator.

**borderSwitch.interfaces [array of interface]**

:   The interfaces belonging to this router to query.

**interface.descr [string]**

:   The description of the interface to use in the output report. Members of the `MonitoredInterfaces`
    array in said report are constructed as: `borderSwitch.hostname_interface.descr`.

**interface.snmpIndex [number]**

:   The index to use when looking into either `ifXTable` or `ifTable`, depending on the configuration
    of `borderSwitch.hcSupport`. This index can be obtained by manually questioning the switch with
    well-known SNMP clients such as `snmpget(1)`.

**outputs [object]**

:   `wlcg-site-snmp` can provide the generated report in several ways. The report can
    be stored 'as-is' on a regular file on disk or it can be served as a regular JSON
    document through both an HTTP and an HTTPS server. All these outputs can be mixed
    and matched depending on the configuration of this option.

**outputs.file [object]**

:   This option defines how (and if) the report is to be stored as a regular file on disk.

**outputs.file.enabled [bool]**

:   Whether to dump the generated report to a file on disk. If `false`, `outputs.file.path`
    will be silently ignored.

**outputs.file.path [string]**

:   The path to dump the generated report on. Bear in mind the user running the program
    must be granted the appropriate permissions to create said file. This file will be
    truncated every `interval` seconds.

**outputs.server [object]**

:   This option defines how (and if) the report is to be served over builtin HTTP and HTTPS servers.
    As long as the ports each of them bind to are different, both servers can be enabled at the same
    time.

**outputs.server.http [object]**

:   This option defines the different settings of the builtin HTTP server.

**outputs.server.http.enabled [bool]**

:   Whether to enable the HTTP server. If false, the rest of `outputs.server.http.*` settings
    are silently ignored.

**outputs.server.http.bindAddress [string]**

:   The address to bind the server to. If set to `0.0.0.0`, the server will listen on every
    available interface. If set to `127.0.0.1`, it will only listen on the local (i.e. `lo`)
    interface.

**outputs.server.http.bindPort [number]**

:   The port to bind the server to. Bear in mind that if this port is lower than `1024` the user
    running the program must be granted the privileges needed to bind to the ports. This can be
    accomplished without having to run the program as `root` through capabilities. Check
    `capabilities(7)` for more information on that.

**outputs.server.https [object]**

:   This option defines the different settings of the builtin HTTPS server.

**outputs.server.https.enabled [bool]**

:   Whether to enable the HTTPS server. If false, the rest of `outputs.server.https.*` settings
    are silently ignored.

**outputs.server.https.bindAddress [string]**

:   The address to bind the server to. If set to `0.0.0.0`, the server will listen on every
    available interface. If set to `127.0.0.1`, it will only listen on the local (i.e. `lo`)
    interface.

**outputs.server.https.bindPort [number]**

:   The port to bind the server to. Bear in mind that if this port is lower than `1024` the user
    running the program must be granted the privileges needed to bind to the ports. This can be
    accomplished without having to run the program as `root` through capabilities. Check
    `capabilities(7)` for more information on that.

**outputs.server.https.certPath [string]**

:   The path of the file containing the server's certificate. Bear in mind the user running
    the program must be capable of reading the provided path. Common formats such as PEM
    certificates should work without any problems.

**outputs.server.https.keyPath [string]**

:   The path of the file containing the server's private key. Bear in mind the user running
    the program must be capable of reading the provided path. Common formats such as PEM
    keys should work without any problems.

**outputs.scp [object]**

:   This option defines how (and if) the report is to be copied to a remote destination over scp(1).
    As scp(1) relies on ssh(1) for transferring data we will refer to several ssh(1) configuration
    parameters. If the output configuration is not working as intended, make sure you can run scp(1)
    with the provided parameters 'as-is' on a shell as the same user that's running the program. If
    everything works as expected, this makes discarding possible permission problems and other system
    misconfigurations that much easier!

**outputs.scp.enabled [bool]**

:   Whether to copy the generated report to a remote machine. If `false`, the following `outputs.scp.*`
    options will be silently ignored.

**outputs.scp.user [string]**

:   The user to to establish the ssh(1) session as. If using ssh(1), this would be the content of the
    connection URI before the **@** sign.

**outputs.scp.hostname [string]**

:   The hostname to ssh(1) into. If using ssh(1), this would be the content of the connection URI
    after the **@** sign.

**outputs.scp.port [number]**

:   The port to establish the connection to `hostname` on. For ssh(1) this is usually **22** unless
    otherwise specified.

**outputs.scp.privateKeyPath [string]**

:   The **absolute path** to the private key to authenticate the ssh(1) session with. This key should
    be readable by the user running the program. Bear in mind the common permission configuration for
    private keys (like those generated with ssh-keygen(1)) is **0600**, which is quite restrictive.

**outputs.scp.serverPublicKey [string]**

:   The remote server's public key. The easiest way to configure this parameter is to manually ssh(1)
    into the server and then look for its hostname on `$HOME/.ssh/known_hosts`. At any rate, this
    value should be the concatenation of the key's type (i.e. `ssh-rsa`, `ecdsa-sha2-nistp256`...)
    with the base64(1)-encoded key. The format itself is specified on section **SSH_KNOWN_HOSTS FILE FORMAT**
    of sshd(8). Even though it is **NOT RECOMMENDED**, this value can be left empty (i.e. `""`) to disable
    key validation. This will, however, emit a warning message on the log each time a file is copied just
    to be extra annoying :P. At any rate, you can automatically find the appropriate value with a simple
    awk(1) invocation. For a remote machine whose hostname is `acme.corp`:

    $ awk awk '/acme\.corp,.+/  {printf("%s %s\n", $2, $3)}' $HOME/.ssh/known_hosts

**outputs.scp.remotePath [string]**

:   The path to copy the file to on the remote machine. This is equivalent to the path specified
    after the colon (i.e. `:`) when invoking scp(1) on a shell.

**outputs.scp.permissions [string]**

:   The numeric mode permissions to apply to the remote file once it's copied as if specified to chmod(1).
    Bear in mind the leading zero can be omitted, but we advise against it. If the desired permissions
    for the file are, for instance, `rw-r--r--` you should specify `0664`.

# AUTHORS
- Pablo Collado Soto <pcolladosoto@gmx.com>
