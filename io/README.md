# Wrappers for golang io

*NOTE*: These aren't exposed to LUA directly, but used as utilities for other classes, which need 
the bridge between lua file and these interfaces: 

- io.Reader
- io.Writer
- io.Seeker
- io.Closer

See usages of `CheckIOReader` and `CheckIOWriter` in json and yaml modules.
