package v2

import (
	"archive/zip"
	"bytes"
	"context"
	"cook-robot-middle-platform-go/config"
	"cook-robot-middle-platform-go/grpc"
	pb "cook-robot-middle-platform-go/grpc/commandRPC"
	"cook-robot-middle-platform-go/httpServer/model"
	"cook-robot-middle-platform-go/logger"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
	"image/png"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type System struct {
	grpcClient *grpc.GRPCClient

	ws            *websocket.Conn
	isUpdating    bool // 是否正在更新软件
	updatingMutex sync.Mutex
}

func NewSystem(grpcClient *grpc.GRPCClient) *System {
	return &System{
		grpcClient: grpcClient,
		isUpdating: false,
	}
}

func (s *System) Shutdown(ctx *gin.Context) {
	req := &pb.ShutdownRequest{
		Empty: true,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, _ := s.grpcClient.Client.Shutdown(ctxGRPC, req)
	logger.Log.Printf("controller关闭成功%d", res)
	os.Exit(1)
}

func (s *System) GetQrCode(ctx *gin.Context) {
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Log.Println("Error:", err)
		model.NewFailResponse(ctx, err)
		return
	}

	// 遍历所有网络接口
	for _, iface := range ifaces {
		// 筛选出WLAN接口，可以根据具体的名称进行判断
		if iface.Name == "wlan0" || iface.Name == "Wi-Fi" || iface.Name == "WLAN" {
			addrs, err := iface.Addrs()
			if err != nil {
				logger.Log.Println("Error:", err)
				model.NewFailResponse(ctx, err)
				return
			}

			// 遍历该接口的IP地址
			for _, addr := range addrs {
				// 检查是否是IPv4地址
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					//logger.Log.Println("WLAN IP Address:", ipnet.IP.String())
					qr, err := qrcode.New("phonePairing::"+ipnet.IP.String()+"\r\n", qrcode.Medium)
					if err != nil {
						logger.Log.Println("Error:", err)
						model.NewFailResponse(ctx, err)
						return
					}

					var buf bytes.Buffer
					png.Encode(&buf, qr.Image(256))

					encodedQrImage := base64.StdEncoding.EncodeToString(buf.Bytes())
					model.NewSuccessResponse(ctx, encodedQrImage)
					return
				}
			}
		}
	}

	model.NewFailResponse(ctx, "no ip found")
}

func (s *System) CheckUpdatePermission(ctx *gin.Context) {
	// 检查控制器是否处在运行状态，运行状态下不允许更新
	if s.grpcClient.ControllerStatus.IsRunning || s.isUpdating {
		model.NewSuccessResponse(ctx, gin.H{
			"isRunning":       s.grpcClient.ControllerStatus.IsRunning,
			"isUpdating":      s.isUpdating,
			"updatePermitted": false,
		})
	} else {
		model.NewSuccessResponse(ctx, gin.H{
			"isRunning":       s.grpcClient.ControllerStatus.IsRunning,
			"isUpdating":      s.isUpdating,
			"updatePermitted": true,
		})
	}
}

func (s *System) Update(ctx *gin.Context) {
	if s.isUpdating {
		logger.Log.Println("正在更新中，拒绝再次更新")
		return
	}
	s.isUpdating = true

	defer func() {
		s.isUpdating = false
	}()

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	logger.Log.Println("建立WebSocket连接")
	s.ws = conn
	defer func() {
		conn.Close()
		s.ws = nil
		logger.Log.Println("断开WebSocket连接")
	}()

	fileURL := fmt.Sprintf("%s:%d/%s", config.App.SoftwareUpdate.ServerHost, config.App.SoftwareUpdate.ServerPort,
		config.App.SoftwareUpdate.Filename)
	err = s.downloadAndSaveFile(fileURL)
	if err != nil {
		logger.Log.Printf("downloadAndSaveFile error:%s", err.Error())
		return
	}

	zipFile := filepath.Join(config.App.SoftwareUpdate.SavePath, config.App.SoftwareUpdate.Filename)
	err = s.unzipFile(zipFile)
	if err != nil {
		logger.Log.Printf("unzipFile error:%s", err.Error())
		return
	}
}

func (s *System) downloadAndSaveFile(fileURL string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		return err
	}

	// 获取文件的总大小
	totalSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	// 创建本地文件
	file, err := os.Create(filepath.Join(config.App.SoftwareUpdate.SavePath, config.App.SoftwareUpdate.Filename))
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 10240) // 缓冲区大小可以根据需求调整
	startTime := time.Now()
	lastTime := startTime
	lastBytes := 0
	totalBytes := 0

	var downloadSpeed float64 = 0
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err = file.Write(buf[:n])
			if err != nil {
				return err
			}
			totalBytes += n

			currentTime := time.Now()
			elapsedTime := currentTime.Sub(lastTime).Milliseconds()
			if elapsedTime > 1000 {
				downloadSpeed = float64(totalBytes-lastBytes) / (float64(elapsedTime) / 1000) / 1024 / 1024 // MB/s
				lastTime = currentTime
				lastBytes = totalBytes + n
			}

			// 实时发送下载进度到前端
			downloadProgress := float64(totalBytes) / float64(totalSize)
			err = s.sendUpdateData(false, false, downloadProgress, 0, downloadSpeed, 0)
			if err != nil {
				return err
			}

		}
		if err == io.EOF {
			logger.Log.Println("下载完毕")
			break
		}
		if err != nil {
			return err
		}
	}
	err = s.sendUpdateData(true, false, 1, 0, 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s *System) unzipFile(zipFile string) error {

	// 打开ZIP文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标文件夹
	if err := os.MkdirAll(config.App.SoftwareUpdate.UnzipPath, 0755); err != nil {
		return err
	}

	totalFiles := len(r.File)
	completedFiles := 0

	removeUIFolderFlag := false

	// 遍历ZIP文件中的每个文件
	for _, file := range r.File {
		// 构建解压后的文件路径
		extractedFilePath := filepath.Join(config.App.SoftwareUpdate.UnzipPath, file.Name)
		// 如果文件是文件夹，创建对应的文件夹
		if file.FileInfo().IsDir() {
			// 如果压缩包中含有electron ui的打包文件夹，则先删除后再解压
			if strings.Contains(file.Name, config.App.SoftwareUpdate.UIFolderName) && !removeUIFolderFlag {
				removeUIFolderFlag = true
				uiFolderPath := filepath.Join(config.App.SoftwareUpdate.SavePath, config.App.SoftwareUpdate.UIFolderName)
				logger.Log.Printf("发现%s文件夹，删除\n", uiFolderPath)
				err = os.RemoveAll(uiFolderPath)
				if err != nil {
					return err
				}
			}
			err = os.MkdirAll(extractedFilePath, file.Mode())
			if err != nil {
				return err
			}
		} else {
			if strings.Contains(file.Name, config.App.SoftwareUpdate.MiddlePlatformFilename) {
				middlePlatformFilePath := filepath.Join(config.App.SoftwareUpdate.SavePath, config.App.SoftwareUpdate.MiddlePlatformFilename)
				logger.Log.Printf("发现%s文件，删除\n", middlePlatformFilePath)
				err = os.RemoveAll(middlePlatformFilePath)
				if err != nil {
					return err
				}
			}

			if strings.Contains(file.Name, config.App.SoftwareUpdate.ControllerFilename) {
				controllerFilePath := filepath.Join(config.App.SoftwareUpdate.SavePath, config.App.SoftwareUpdate.ControllerFilename)
				logger.Log.Printf("发现%s文件，删除\n", controllerFilePath)
				err = os.RemoveAll(controllerFilePath)
				if err != nil {
					return err
				}
			}

			// 否则，创建上层文件夹并解压文件
			if err = os.MkdirAll(filepath.Dir(extractedFilePath), 0755); err != nil {
				return err
			}
			// 打开ZIP文件中的文件
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			// 创建目标文件
			dstFile, err := os.Create(extractedFilePath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			// 将ZIP文件中的内容复制到目标文件
			_, err = io.Copy(dstFile, rc)
			if err != nil {
				return err
			}

			completedFiles++
			unzipProgress := float64(completedFiles) / float64(totalFiles)
			err = s.sendUpdateData(true, false, 1, unzipProgress, 0, 0)
			if err != nil {
				return err
			}
		}
	}
	logger.Log.Println("解压完毕")

	err = s.sendUpdateData(true, true, 1, 1, 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s *System) sendUpdateData(isDownloadFinished bool, isUnzipFinished bool,
	downloadProgress float64, unzipProgress float64, downloadSpeed float64, unzipSpeed float64) error {
	//logger.Log.Println(downloadProgress, unzipProgress, downloadSpeed, unzipSpeed)
	err := s.ws.WriteJSON(gin.H{
		"isDownloadFinished": isDownloadFinished,
		"isUnzipFinished":    isUnzipFinished,
		"downloadProgress":   downloadProgress,
		"unzipProgress":      unzipProgress,
		"downloadSpeed":      downloadSpeed,
		"unzipSpeed":         unzipSpeed,
	})
	if err != nil {
		logger.Log.Printf("error:%s", err.Error())
	}
	return err
}
