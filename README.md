sft
===

## Simple file transfer utility

Fastest file transfer on Linux.
No encription, simple TCP connection.

10x faster than transfer over SSH link.
5x faster than NetCat connection.

### Examples

#### From client to server
Server:

		$ ./sft --destination file.out
Client:

		$ ./sft 127.0.0.1 file.in

#### From server to client
Server:

		$ ./sft --source file.in
Client:

		$ ./sft 127.0.0.1 file.out

#### With pipes:
Server:

		$ ./sft --source /dev/zero
Client:

		$ ./sft --verbose 127.0.0.1 - | pv > /dev/null
		Client mode
		Port: 18000
		File: -
		Connecting... ok
		Mode: Source
		4.37GB 0:00:05 [ 909MB/s] [    <=>                                             ]

Copyright (c) 2014 Dmitry Lagoza
