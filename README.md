每隔3-5分钟随机休息机制实现。

默认周期90分钟，休息20分钟。

NOTE: build不提供mp3文件，需要自己添加mp3文件到mp3目录: start.mp3、replay.mp3、finish.mp3

```sh
make

./random-replay.exe --help
./random-replay.exe -p 90 -b 20
```
