package app

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/jebo87/makako-gateway/clients"
	"gitlab.com/jebo87/makako-grpc/ads"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	router = gin.Default()
	kacp   = keepalive.ClientParameters{
		Time:                3 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,     // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,            // send pings even without active streams
	}
)

func StartApp() {
	MapURLs()
	grpcConnection := connectGRPC()

	defer grpcConnection.Close()
	clients.GrpcClient = ads.NewAdsClient(grpcConnection)
	router.Run("localhost:8080")
}

func connectGRPC() *grpc.ClientConn {
	conn, err := grpc.Dial(os.Getenv("API_ADDRESS_DEV")+":"+os.Getenv("API_PORT_DEV"), grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))

	log.Println("connecting to GRPC server " + os.Getenv("API_ADDRESS_DEV"))

	if err != nil {
		panic(err) // TODO: define a way to keep trying without crashing the app
	}
	return conn
}
