package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	CHILD_ENV_KEY = "GO_CHILD_ENV_KEY"
	CHILD_ENV_VAL = "GO_CHILD_VALUE"
)

//const (
//	LOG1_FILE = "/var/log/go/parent.log"
//	LOG2_FILE = "/var/log/go/child.log"
//)

func main() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	cmd, err := NewProcess(false)
	if err != nil {
		log.Printf("new process err: %s \n", err.Error())
		return
	}

	if cmd == nil { // child
		childProcess()
		return
	}

	// parent
	log.Printf("parent pid: %d \n", os.Getpid())
	quit := make(chan struct{})
	restartCh := make(chan struct{})
	pidCh := make(chan int, 1)

	// 监听事件
	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			event := rand.Intn(100)
			if event == 0 {
				// 发送重启信号
				log.Println("send restart")
				select {
				case restartCh <- struct{}{}:
				case <-quit:
					return
				}
			}
		}
	}()

	// 监听退出信号
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		s := <-signals

		log.Println("parent receive signal: " + s.String())

		// 广播退出信号
		close(quit)
	}()

	// 监听重启事件和退出信号发送给子进程
	go func() {
		for {
			select {
			case <-restartCh:
			case <-quit:
			}

			select {
			case pid := <-pidCh:
				// 收到子进程pid
				_ = syscall.Kill(pid, syscall.SIGUSR1)
			default:
				// 没有收到时，不做操作，避免短时间内多次操作
			}
		}
	}()

	for {
		pidCh <- cmd.Process.Pid

		if err := cmd.Wait(); err != nil {
			log.Printf("process child err: %s \n", err.Error())
			// 子进程错误退出，父进程也要退出
			break
		}

		select {
		case <-quit:
			// 退出
			return
		default:
			cmd, err = NewProcess(false)
			if err != nil {
				log.Printf("restart process err: %s \n", err.Error())
				return
			}
		}
	}
}

func childProcess() {
	log.Printf("child pid: %d \n", os.Getpid())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	s := <-signals
	log.Println("child exit, signal: " + s.String())
}

func NewProcess(parentExit bool) (*exec.Cmd, error) {
	val := os.Getenv(CHILD_ENV_KEY)
	if val == CHILD_ENV_VAL {
		// child process
		return nil, nil
	}

	//pf, err := os.Create(LOG1_FILE)
	//if err != nil {
	//	return nil, err
	//}
	//log.SetOutput(pf)
	//
	//f, err := os.Create(LOG2_FILE)
	//if err != nil {
	//	return nil, err
	//}
	//defer f.Close()

	childEnv := append(os.Environ(), fmt.Sprintf("%s=%s", CHILD_ENV_KEY, CHILD_ENV_VAL))
	cmd := exec.Cmd{
		Path:   os.Args[0],
		Args:   os.Args,
		Env:    childEnv,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start child process: %w", err)
	}

	if parentExit {
		os.Exit(0)
	}

	return &cmd, nil
}
