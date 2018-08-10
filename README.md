# pipe-tar

## Usage
1. Examples
    * find /tmp/ -type f | ./pipe-tar --stdin --abs-path --output-file tmp.tar
    * find /tmp/ -type f | tar -cvf tmp.tar -T-
2. Options
    * --stdin: get file list from stdin
    * --delim: delimiter, default is '\n'
    * --path: file path
    * --format: Create archive of the given format.
        ```
        gnu                      GNU tar 1.13.x format
        oldgnu                   GNU format as per tar <= 1.12
        pax                      POSIX 1003.1-2001 (pax) format
        posix                    same as pax
        ustar                    POSIX 1003.1-1988 (ustar) format
        v7                       old V7 tar format
        ```
    * --compress:
    * --abs-path: abs path
    * --output-file: output file
