package rpc

// Ping ...
import (
	"github.com/c3systems/c3-go/config"
	"github.com/c3systems/c3-go/node/store/disk"
	pb "github.com/c3systems/c3-go/rpc/pb"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func ping() *pb.LatestBlockResponse {
	cnf := config.New()
	//diskStore, err := leveldbstore.New(cfg.DataDir(), nil)
	diskStore, err := disk.New(&disk.Props{})
	if err != nil {
		log.Fatal(err)
	}

	headBlock, err := diskStore.GetHeadBlock()
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(headBlock)

	return &pb.LastBlockResponse{
		Data: "block",
	}
}
