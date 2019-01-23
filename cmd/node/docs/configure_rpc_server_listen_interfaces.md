pod allows you to bind the RPC server to specific interfaces which enables you to setup configurations with varying levels of complexity. The `rpclisten` parameter can be specified on the command line as shown below with the -- prefix or in the configuration file without the -- prefix (as can all long command line options).

The configuration file takes one entry per line.

A few things to note regarding the RPC server:

- The RPC server will **not** be enabled unless the `rpcuser` and `rpcpass` options are specified.

- When the `rpcuser` and `rpcpass` and/or `rpclimituser` and `rpclimitpass` options are specified, the RPC server will only listen on localhost IPv4 and IPv6 interfaces by default. You will need to override the RPC listen interfaces to include external interfaces if you want to connect from a remote machine.

- The RPC server has TLS disabled by default. You may use the `--TLS` option to enable it.

- The `--rpclisten` flag can be specified multiple times to listen on multiple interfaces as a couple of the examples below illustrate.

- The RPC server is disabled by default when using the `--regtest` and
  `--simnet` networks. You can override this by specifying listen interfaces.

Command Line Examples:

| Flags                                           | Comment                                                             |
| ----------------------------------------------- | ------------------------------------------------------------------- |
| --rpclisten=                                    | all interfaces on default port which is changed by `--testnet`      |
| --rpclisten=0.0.0.0                             | all IPv4 interfaces on default port which is changed by `--testnet` |
| --rpclisten=::                                  | all IPv6 interfaces on default port which is changed by `--testnet` |
| --rpclisten=:11048                              | all interfaces on port 11048                                        |
| --rpclisten=0.0.0.0:11048                       | all IPv4 interfaces on port 11048                                   |
| --rpclisten=[::]:11048                          | all IPv6 interfaces on port 11048                                   |
| --rpclisten=127.0.0.1:11048                     | only IPv4 localhost on port 11048                                   |
| --rpclisten=[::1]:11048                         | only IPv6 localhost on port 11048                                   |
| --rpclisten=:8336                               | all interfaces on non-standard port 8336                            |
| --rpclisten=0.0.0.0:8336                        | all IPv4 interfaces on non-standard port 8336                       |
| --rpclisten=[::]:8336                           | all IPv6 interfaces on non-standard port 8336                       |
| --rpclisten=127.0.0.1:8337 --listen=[::1]:11048 | IPv4 localhost on port 8337 and IPv6 localhost on port 11048        |
| --rpclisten=:11048 --listen=:8337               | all interfaces on ports 11048 and 8337                              |

The following config file would configure the pod RPC server to listen to all interfaces on the default port, including external interfaces, for both IPv4 and IPv6:

```text
[Application Options]

rpclisten=
```

As well as the standard port 11048, which delivers sha256d block templates, there is another 8 RPC ports that open sequentially up to 11056 with each one providing a specific block template version number.