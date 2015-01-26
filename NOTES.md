
btcwallet changes
=================

In rpcserver.go L: 1300

```go
// TODO NOTICE
// This is a custom RPC command that handles sending bulletins
"sendbulletin": SendBulletin,
```

In rpcserver.go L:40

```go
// TODO NOTICE import package derp for effect
_ "github.com/soapboxsys/ombudslib/walletexten"
```


btcd changes
============

On ishtar in NSkelsey fork


btcctl changes
=============

In btcctl.go L:13
```go
    // TODO NOTICE
    _"github.com/soapboxsys/ombudslib/walletexten"
```

In btcctl.go L:115
```go
// TODO NOTICE
// Added command, a conversion handler just moves the types
"sendbulletin": {3, 0, displayGeneric, []conversionHandler{nil, nil, nil}, makeSendBulletin, "<address> <board> <message>"},
```

In btcctl.go L:223
```go
// TODO NOTICE
// makes a sendBulletin cmd from the provided parameters
func makeSendBulletin(args []interface{}) (btcjson.Cmd, error) {
        return walletexten.NewSendBulletinCmd("btcctl", args[0].(string), args[1].(string), args[2].(string)), nil
}
```

