# Gombadi Reflector

This repo contains an http application that will reflect back what you send to it.

Set it up on a server and visit the page.

http:your.domainname.com/ - This will return a standard message

http:your.domainname.com/ip - This will return your ip address

http:your.domainname.com/all - This will return all details of the request made.

If you ask for /all then you can also add ?o=json or ?o=xml to see the output in the requested format.



## Installing

Simply use go get to download the code:

    $ go get github.com/gombadi/http-reflector


Then run it using ./reflector -p <8880>

If you then connect to it with a browser you will see the resquest details that were
sent through to the server.

