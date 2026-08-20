### SolarmanV5 client & library

#### Client

High-level client for easy access to solar/hybrid inverters data through IGEN/Solarman dataloggers.

#### Packages

The V5 protocol and MODBUS packages can be used as standalone libraries should you with to create your own client

*Note: The `modbus` package is not a complete MODBUS protocol implementation*


#### Install
```
go get github.com/githubDante/solarman
```

### Protocol documentation

For an excelent in-depth breakdown of the SolarmanV5 protocol, please refer to the documentation provided by [@jmccrohan](https://github.com/jmccrohan) [here](https://pysolarmanv5.readthedocs.io/en/stable/solarmanv5_protocol.html)

#### Tests

* The MODBUS validation tests use data from the examples found [here](https://simplymodbus.ca/learn-rtu.html)


### Examples
#### Client
* [Read only](examples/solarman-client) client.
* [Verbose read/write](examples/solarman-verbose-client) client with the logger from `github.com/githubDante/solarman/log` package
