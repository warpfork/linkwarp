Add go version to catalog
`warpforge catalog add tar go.dev/dl/go:1.18.7:linux-amd64 https://go.dev/dl/go1.18.7.linux-amd64.tar.gz`

Stupid error info:
```
	error: serialization error: unable to deserialize plot: Invalid byte while expecting start of value: 0x2f
```
plot.wf
```
00000000: 7b0a 0922 706c 6f74 2e76 3122 3a20 7b0a  {.."plot.v1": {.
00000010: 0909 2269 6e70 7574 7322 3a20 7b0a 0909  .."inputs": {...
00000020: 0922 726f 6f74 6673 223a 2022 6361 7461  ."rootfs": "cata
00000030: 6c6f 673a 7761 7270 7379 732e 6f72 672f  log:warpsys.org/
00000040: 6275 7379 626f 783a 7631 2e33 352e 303a  busybox:v1.35.0:
```
crash
```
$ warpforge ferk --plot plot.wf --cmd /pkg/busybox/bin/sh
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0xd87731]

goroutine 1 [running]:
main.cmdFerk(0xc0004f02c0)
	/home/cjb/repos/warpforge/wf/cmd/warpforge/ferk.go:83 +0x4f1
github.com/warpfork/warpforge/cmd/warpforge/internal/util.CmdMiddlewareTracingSpan.func1(0xc0004f02c0)
	/home/cjb/repos/warpforge/wf/cmd/warpforge/internal/util/middleware.go:42 +0x138
github.com/warpfork/warpforge/cmd/warpforge/internal/util.CmdMiddlewareTracingConfig.func1(0xc0004f02c0)
	/home/cjb/repos/warpforge/wf/cmd/warpforge/internal/util/middleware.go:62 +0x203
github.com/warpfork/warpforge/cmd/warpforge/internal/util.CmdMiddlewareLogging.func1(0xc0004f02c0)
	/home/cjb/repos/warpforge/wf/cmd/warpforge/internal/util/middleware.go:31 +0x182
github.com/urfave/cli/v2.(*Command).Run(0x17cdf40, 0xc0004f0040)
	/home/cjb/gopath/pkg/mod/github.com/urfave/cli/v2@v2.3.0/command.go:163 +0x5bb
github.com/urfave/cli/v2.(*App).RunContext(0xc0004e4340, {0x11856d8?, 0xc000036110}, {0xc0000321e0, 0x6, 0x6})
	/home/cjb/gopath/pkg/mod/github.com/urfave/cli/v2@v2.3.0/app.go:313 +0xb48
github.com/urfave/cli/v2.(*App).Run(...)
	/home/cjb/gopath/pkg/mod/github.com/urfave/cli/v2@v2.3.0/app.go:224
main.main()
	/home/cjb/repos/warpforge/wf/cmd/warpforge/main.go:124 +0x67
```
