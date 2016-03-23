
ombwebrelay
==========

A RESTful Web API that exposes the public record to the Internet.
The functionality of this api is defined in [the Ombuds spec](spec.getombuds.org/web-api).

Installing
----------

This binary depends on a pubrecord.db that is created and updated by ombfullnode. 
To run this rest server you must download it, build it, and give it the path to the db.
OR you can run the bash script!

```bash
> git clone https://github.com/soapboxsys/ombudslib
> cd ombudslib/cmds/ombwebrelay
> # This next step assumes you have a pubrecord.db at $DBPATH
> ./start.sh
```

Becoming a Relay Operator
-------------------------

Running a relay is the best way you can support the network.
