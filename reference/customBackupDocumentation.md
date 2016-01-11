cPanel custom backup destination documentation
==============================================

!Note: For cPanel & WHM 11.48

!Info: This page is copied from [here](https://documentation.cpanel.net/display/ALD/How+to+Create+a+Custom+Destination) on `2015-04-29`

## Overview ##

!Warning: This feature is for advanced users.

The Backup Configuration feature allows users to create a Custom Destination for their backups. Users can create and use a custom destination instead of the currently available destinations in WHM (FTP, SFTP, WebDAV, or a local directory).

## Script ##

This is the path to a script you provide which implements the custom transport.

### Script commands ###

The script must implement the following commands:

`Command` | `Parameters`
----------|--------------
`get` 		|	`remote file` `local file`
`put` 		| `local file` `remote file`
`ls`			| `Path`
`mkdir`		| `Path`
`chdir`		| `Path`
`rmdir`		| `Path`
`delete`	| `Path`

Backups run each of these commands while the backup file is transported and while the destination is validated.

## Script operation ##

* The script will run once per command.
* The script cannot save state information between commands.
* The connection is not reused between commands. Instead, each time the script runs, the connection to the remote custom destination will be created, and dropped after the script is run.
* Information will be passed into the script via the command line in the following order:
 * Command name
 * Current directory
 * Command specific parameters
 * Host
 * Username

By default, the script is told which directory on the remote system to use for the operation. Since the connection to the directory is dropped between operations, the caller must save the directory used.

If the script writes to STDERR, it will fail. Anything written to STDERR will be logged as part of the failure.

Two commands return data back to the caller by writing to STDOUT:

* `chdir` — Prints the new working directory on the remote system as a result of running the chdir command. This is necessary because the remote system may, based on the contents of the path parameter (special characters, etc.), end up in a different working directory than the caller is able to anticipate.
* `ls` — Will print output identical to what running ls -l.
 * For example: `-rwxr-xr-1 root root 3171 Jan 18 12:23 temp.txt`
!Note: If you entered a password when you created the custom destination, it is passed through to the script via the PASSWORD environment variable. The password will not be displayed as part of the command line.
An example custom destination can be found in `/usr/local/cpanel/scripts/custom_backup_destination.pl.skeleton`
