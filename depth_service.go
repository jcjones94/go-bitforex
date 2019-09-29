package bitforex

import (
	"context"
)

// DepthService show depth info
type DepthService struct {
	c      *Client
	symbol string
	size  *int
}

// Symbol set symbol
func (s *DepthService) Symbol(symbol string) *DepthService {
	s.symbol = symbol
	return s
}

// Limit set dpeth size
func (s *DepthService) Size(size int) *DepthService {
	s.size = &size
	return s
}

// Do send request
func (s *DepthService) Do(ctx context.Context, opts ...RequestOption) (res *DepthResponse, err error) {
	r := &request{
		method:   "GET",
		endpoint: "/api/v1/market/depth",
	}
	r.setParam("symbol", s.symbol)
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	j, err := newJSON(data)
	if err != nil || j.Get("success").MustBool() == false{
		return nil, err
	}

	res = new(DepthResponse)
	res.LastUpdateID = j.Get("time").MustInt64()
        j_data := j.Get("data")
	bidsLen := len(j_data.Get("bids").MustArray())
	res.Bids = make([]Bid, bidsLen)
	for i := 0; i < bidsLen; i++ {
		item := j_data.Get("bids").GetIndex(i)
		res.Bids[i] = Bid{
			Price:    item.Get("price").MustFloat64(),
			Quantity: item.Get("amount").MustFloat64(),
		}
	}
	asksLen := len(j_data.Get("asks").MustArray())
	res.Asks = make([]Ask, asksLen)
	for i := 0; i < asksLen; i++ {
		item := j_data.Get("asks").GetIndex(i)
		res.Asks[i] = Ask{
			Price:    item.Get("price").MustFloat64(),
			Quantity: item.Get("amount").MustFloat64(),
		}
	}
	return res, nil
}

// DepthResponse define depth info with bids and asks
type DepthResponse struct {
	LastUpdateID int64 `json:"lastUpdateId"`
	Bids         []Bid `json:"bids"`
	Asks         []Ask `json:"asks"`
}

// Bid define bid info with price and quantity
type Bid struct {
	Price    float64
	Quantity float64
}

// Ask define ask info with price and quantity
type Ask struct {
	Price    float64
	Quantity float64
}
