package main

import (
	"bytes"
	"context"
	"embed"
	"io"
	"path"

	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/spf13/pflag"
)

var periodTime int
var breakTime int

//go:embed mp3/*.mp3
var mp3Files embed.FS

func init() {
	pflag.IntVarP(&periodTime, "period", "p", 60, "define period time.(per min)")
	pflag.IntVarP(&breakTime, "break", "b", 20, "define break time.(per min)")
}

func main() {
	pflag.Parse()

	if periodTime < 0 || breakTime < 0 {
		fmt.Println("参数不允许取当前值")
		os.Exit(1)
	}
	for {
		fmt.Println("周期开始")
		// 周期开始提示音
		err := playBeep("start.mp3")
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		// 启动协程
		ctx, cancel := context.WithCancel(context.Background())
		go randomReplay(ctx)

		// 等待周期完成
		log.Printf("还有%d分钟", periodTime)
		time.Sleep(time.Minute * time.Duration(periodTime))

		// 结束协程
		fmt.Println("发出终止信号")
		cancel()

		// 周期结束提示音
		fmt.Println("周期结束")
		err = playBeep("finish.mp3")
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		// 等待休息完成
		log.Printf("还有%d分钟", breakTime)
		time.Sleep(time.Minute * time.Duration(breakTime))
		fmt.Println("休息完毕")
	}
}

func playBeep(fs string) error {
	fs = path.Join("mp3", fs)

	data, err := mp3Files.ReadFile(fs)
	if err != nil {
		return fmt.Errorf("open beep file error: %s", err)
	}

	streamer, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return fmt.Errorf("decode error: %s", err)
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return fmt.Errorf("init speaker error: %s", err)
	}
	done := make(chan bool)

	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   -2, // 分贝（负数为降低音量）
		Silent:   false,
	}

	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- true
	})))
	<-done

	return nil
}

func randomReplay(ctx context.Context) {
	for {
		randomSecond := rand.Intn(121) + 180 // 3-5分钟随机，单位秒

		// 每次sleep一秒，监听是否被终止
		for range randomSecond {
			select {
			case <-ctx.Done():
				fmt.Println("协程终止")
				return
			default:
				time.Sleep(time.Second)
			}
		}

		// 播放休息提示音
		fmt.Println("休息十秒钟")
		err := playBeep("replay.mp3")
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		for range 10 {
			select {
			case <-ctx.Done():
				fmt.Println("协程终止")
				return
			default:
				time.Sleep(time.Second)
			}
		}

		// 播放休息结束提示音
		fmt.Println("结束十秒休息")
		err = playBeep("start.mp3")
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
	}
}
