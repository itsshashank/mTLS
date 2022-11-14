package timer

import (
	"context"
	"io"
	"log"

	timeserver "github.com/itsshashank/gtimenow/model"
)

type TimeServer struct {
	timeserver.UnimplementedTimeServerServer
}

func (c *TimeServer) TimeNow(ts timeserver.TimeServer_TimeNowServer) error {
	log.Println("start new server")
	var max int64
	ctx := ts.Context()

	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := ts.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		// continue if number reveived from stream
		// less than max
		r := req.GetWord()
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}
		if r <= max {
			continue
		}

		// update max and send it to stream
		max = r
		resp := timeserver.Response{Word: max}
		if err := ts.SendMsg(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		log.Printf("send new max=%d", max)
	}
}

func (c *TimeServer) Hello(ctx context.Context, in *timeserver.Response) (*timeserver.Response, error) {
	return &timeserver.Response{Word: in.GetWord()}, nil
}
