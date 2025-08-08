# neteasePlaylistAPI
neteasePlaylistAPI 利用NeteaseCloudMusicApi解析playlist并输出歌词/歌曲  
搭配home项目使用

只支持游客模式能够播放的音乐,使用登录模式请自行修改项目文件 ( *不负责教学 请不要为此发送issue或联系我* )  
关于利用本项目进行共享账户导致的封禁与版权问题本项目不负任何责任  
建议将Cookie通过外部文件读入后载入在 api/global.go 内的 `fetchAPI` 函数 (可自行写一个刷新Cookie的程序+部署定时任务执行更新该文件)  

# API
感谢 api.sayqz.com 提供NeteaseCloudMusicApi

# Demo
> [!WARNING]  
> 稳定服务建议完整API都自行搭建而不是使用演示API  
> 如果你要使用演示API用于生产环境请合理使用    
> 如果使用过于频繁可能会被我们封禁  

> API ``dev.moguq.top``  
> 使用演示 [蘑菇の音乐播放器](https://www.moguq.top/music)


# 免责声明
本项目只作为前置API处理透过代理的官方API返回的数据  
如果贵司依旧认为本项目侵犯了您的权益请与我联系  

# License
MIT License
