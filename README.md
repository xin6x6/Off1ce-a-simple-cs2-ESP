# Off1ce
A reaaaaaly simple cs2 esp cheat     使用go语言

<h4>这是我之前的一个叫做espcel.exe项目的改进</h4>

## 编译 compile
  1. 从Go的官网下载并安装Go
      ```
      https://go.dev/doc/install
      ```
  2. cd至源码文件夹
      
  3. 在源码文件夹目录下打开终端，输入:
      ```
      set GOOS=windows
      set GOARCH=amd64
      ```
  4. 输入:
      ```
      go build -ldflags "-s -w"
      ```

## 更新offsets.json
  你可以从cs2-dumper获得最新的`offsets.json`
