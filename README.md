# neteasePlaylistAPI
neteasePlaylistAPI 利用NeteaseCloudMusicApi解析playlist并输出歌词/歌曲  
搭配home项目使用

只支持游客模式能够播放的音乐,使用登录模式请自行修改项目文件 (不负责教学) 
建议将Cookie通过外部文件读入后载入在 api/global.go 内的 `fetchAPI` 函数 (可自行写一个刷新Cookie的程序+部署定时任务执行更新该文件只有static目录下的html文件会被主动暴露)  

# API
感谢 api.sayqz.com 提供NeteaseCloudMusicApi

# Demo
> API dev.moguq.top  
> [蘑菇の音乐播放器](https://www.moguq.top/music)

# License
MIT License
