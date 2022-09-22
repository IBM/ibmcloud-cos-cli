package aspera

import (
	"io"

	sdk "github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
	"gopkg.in/cheggaaa/pb.v1"
)

type Subscriber interface {
	Queued(resp *sdk.TransferResponse)
	Running(resp *sdk.TransferResponse)
	Done(resp *sdk.TransferResponse)
}

type DefaultSubscriber struct{}

// Do Nothing by default
func (b *DefaultSubscriber) Queued(resp *sdk.TransferResponse)  {}
func (b *DefaultSubscriber) Running(resp *sdk.TransferResponse) {}
func (b *DefaultSubscriber) Done(resp *sdk.TransferResponse)    {}

// display progress bar for file transfer
type ProgressBarSubscriber struct {
	bar *pb.ProgressBar
}

func NewProgressBarSubscriber(total int64, out io.Writer) *ProgressBarSubscriber {
	bar := pb.New64(total).SetUnits(pb.U_BYTES)
	bar.Output = out
	return &ProgressBarSubscriber{bar: bar}
}

func (p *ProgressBarSubscriber) Queued(resp *sdk.TransferResponse) {
	p.bar.Prefix("Queued")
	p.bar.Start()
}

func (p *ProgressBarSubscriber) Running(resp *sdk.TransferResponse) {
	p.bar.Prefix("Running")
	p.bar.Set(int(resp.TransferInfo.BytesTransferred))
}

func (p *ProgressBarSubscriber) Done(resp *sdk.TransferResponse) {
	p.bar.Finish()
}
