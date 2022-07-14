# No Internet Secured Fixer

![GitHub](https://img.shields.io/github/license/lgiki/no-internet-secured-fixer?style=flat-square)

[English](README.md)



当网络连接正常时（可以正常上网），Windows 经常显示 `无 Internet，安全` 。有一些应用程序依靠 Windows 报告的网络状态来确定网络连接是否可用，例如 Spotify ，当 Windows 报告 "无 Internet，安全" 时，它不让我播放任何歌曲。同时，如果 Windows 认为没有网络连接，我就不能打开热点，这非常烦人。

这个问题的主要是由 Windows 的 [Network Connectivity Status Indicator (NCSI)](https://docs.microsoft.com/en-us/troubleshoot/windows-client/networking/internet-explorer-edge-open-connect-corporate-public-network) 导致的。默认情况下，NCSI 会向微软的服务器发送 HTTP 请求和 DNS 请求，并检查返回的结果是否与注册表中保存的结果相符，如果相符，则认为网络正常。

然而，由于服务器不稳定或其他复杂的原因，这个检查过程可能会失败。当NCSI不能得到正确的响应时，它就会认为没有网络连接，导致Windows显示 "无 Internet，安全" 。本程序提供了一些替代的服务器来替代 NCSI 默认使用的微软服务器，你可以在 [servers.json](servers.json) 文件中查看本程序使用的所有服务器。

这个程序无法帮助你解决网络连接的错误，这个程序只能帮你解决网络连接正常时 Windows 显示 "无Internet，安全" 的问题。

因为修改 NCSI 设置需要修改系统注册表，所以这个程序需要管理员权限。

#  截图

![](screenshot.png)

# 用法

- 从 [Releases](https://github.com/LGiki/no-internet-secured-fixer/releases) 页面下载该程序，双击打开并赋予管理员权限。
- 按照菜单提示，选择你需要的操作。通常，你只需要选择第一项（Set NCSI registries automatically），程序就会测试 [servers.json](servers.json) 文件中的所有服务器，并自动选择延迟最小的服务器作为你的NCSI服务器。
- 在控制面板里禁用并重新启用你的网络连接或者重新启动计算机。

# 参考

- [https://github.com/crazy-max/WindowsSpyBlocker](https://github.com/crazy-max/WindowsSpyBlocker)

# 开源协议

MIT