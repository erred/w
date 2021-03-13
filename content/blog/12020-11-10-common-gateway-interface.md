---
description: notes on cgi things
title: common gateway interface
---

### _common_ gateway interface

So how do programs of yore run as web servers
when they didn't include a http stack?
You run a dedicated server process (Apache, Nginx, ...)
which calls your program to generate the response.

#### _CGI_

[rfc3875](https://tools.ietf.org/html/rfc3875)

The oldest protocol.
The server executes a program,
passing request info as env and request body as stdin,
expecting a response (including headers) on stdout.

#### _FCGI_

[spec](http://www.mit.edu/~yandros/doc/specs/fcgi-spec.html#:~:text=1.-,Introduction,Web%20server%20that%20supports%20FastCGI.)

FastCGI, because fork-exec for every request is too costly.
Basically a worker pool of processes that are kept alive,
communication (request,response) is over sockets

#### _SCGI_

[spec](https://python.ca/scgi/protocol.txt)

SimpleCGI, like FCGI but simpler(?) to parse

#### _uwsgi_

[spec](https://uwsgi-docs.readthedocs.io/en/latest/Protocol.html)

Not to be confused with the framework of the same name,
a simplish binary packet protocol for between the server and framework.

#### _JServ_

[spec](http://tomcat.apache.org/tomcat-3.3-doc/AJPv13.html)

Binary packet protocol between web server and application thing.

#### _Language_ Specific

There are others, like for Perl, Clojure, Common Lisp, ...

##### _JavaScript_ JSGI

javascript object as request, javascript object as response.

##### _Python_ WSGI

Web Server Gateway Interface, Python specific.
The python function signature:
request as env/dict, response through callback/return string.
Uses a separate middleware process to manage lifetimes
and translate from HTTP/other to the function call.
Example: `Apache <-HTTP-> uwsgi <-function call-> app`
