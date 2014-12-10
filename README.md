agent-fs
========

Watches a file system tree at root. Whenever a change is noticed, 
parses the tree into a QLM structure and sends it to the core.

Usage: ```agent-fs <root directory> <core address>```

Core address is a websocket and must end in ```/qlmws``` eg. ```ws://host.example.com/qlmws```.

Fs vs. QLM mapping
------------------

Generally QLM elements that can contain other elements or multiple values are represented as directories and elements that only have a single value are represented as files.

Root directory must contain a directory named ```Objects```.
The ```Objects``` directory may contain any number of ```Object.``` directories.
The ```Objects``` directory may also contain following files: ```XmlnsXsi```, ```NoNamespaceSchemaLocation``` and ```Version```. Contents of the files will be used as the values of the corresponding QLM elements.

Object nodes are represented by directories whose name start in ```Object.```.
An ```Object.``` directory may contain any number of ```Object.``` and ```InfoItem.``` directories as well as a directory named ```Id``` and a directory named ```Description```. An ```Object.``` directory may also contain following files: ```Type``` and ```Udef```. Contents of the files will be used as the values of the corresponding QLM elements.

InfoItem nodes are represented by directories whose name start in ```InfoItem.```.
An ```InfoItem.``` directory may contain any number of ```Value.``` directories, a single ```Descripton``` directory and a single ```MetaData``` directory as well as ```Udef```, ```Name``` and ```OtherNames``` files. Contents of the files will be used as the values of the corresponding QLM elements. ```OtherNames``` may contain several names each in its own line.

MetaData nodes are represented by directories names ```MetaData```.
A ```MetaData``` directory may contain any number of ```InfoItem.``` directories.

Values are represented by directories whose name start in ```Value.```.
A ```Value.``` directory may contain files named ```Text```, ```Type```, ```DateTime``` and ```UnixTime```. Contents of the files will be used as the values of the corresponding QLM elements.

Ids are represented by directories named ```Id```.
An ```Id``` directory may contain files named ```IdType```, ```TagType```, ```StartDate```, ```EndDate```, ```Udef``` and ```Text```. Contents of the files will be used as the values of the corresponding QLM elements.

Descriptions are represented by directories named ```Description```.
A ```Description``` directory may contain files named ```Lang```, ```Udef``` and ```Text```. Contents of the files will be used as the values of the corresponding QLM elements.

In the  ```Object.```, ```InfoItem.``` and ```Value.``` directory names the part following the dot is to keep directory names unique and may be anything.

Generating test data
--------------------

The files ```agent-fs-handler.py``` and ```agent_fs.ini``` are provided for use with ```https://github.com/burmanm/gather_agent```. Running ```gather_agent``` with the ini file will use the handler to generate a directory tree that agent-fs can read and to periodically update the values contained in it.
