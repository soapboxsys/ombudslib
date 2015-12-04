### Protocol Buffers
We use google protocol buffers to encode data in output scripts of bitcoin transactions.

This gives us a language agnostic specification

To contribute to this project or build tools for it. You need a protocol buffer extension
for the langauge you are operating in.

For golang the compiler was retreived from [here](http://code.google.com/p/goprotobuf/). 
The file types.pb.go was built using this command:
```bash
> protoc --go_out=./ types.proto
```

If you want to build all of the langauges supported then run:
```bash
> make
```

