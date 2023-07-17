package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"homework9/internal/model/ads"
	"homework9/internal/model/filter"
	"time"
)

// AdToAdResponse converts ads.Ad to *AdResponse
func AdToAdResponse(ad ads.Ad) *AdResponse {
	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
		CreationDate: &Date{
			Day:   int64(ad.CreationDate.Day),
			Month: int64(ad.CreationDate.Month),
			Year:  int64(ad.CreationDate.Year),
		},
		UpdateDate: &Date{
			Day:   int64(ad.UpdateDate.Day),
			Month: int64(ad.UpdateDate.Day),
			Year:  int64(ad.UpdateDate.Year),
		},
	}
}

// AdsListToListAdsResponse converts []ads.Ad to * ListAdResponse
func AdsListToListAdsResponse(listAd []ads.Ad) *ListAdResponse {
	var resList ListAdResponse
	resList.List = make([]*AdResponse, 0, len(listAd))
	for _, a := range listAd {
		resList.List = append(resList.List, AdToAdResponse(a))
	}
	return &resList
}

func (s *Server) CreateAd(ctx context.Context, req *CreateAdRequest) (*AdResponse, error) {
	createdAd, createErr := s.App.CreateAd(ctx, req.Title, req.Text, req.AuthorId)
	if createErr != nil {
		return nil, ErrorToGRPCError(createErr)
	}

	return AdToAdResponse(createdAd), nil
}

func (s *Server) ChangeAdStatus(ctx context.Context, req *ChangeAdStatusRequest) (*AdResponse, error) {
	changedAd, changeErr := s.App.ChangeAdStatus(ctx, req.Id, req.AuthorId, req.Published)
	if changeErr != nil {
		return nil, ErrorToGRPCError(changeErr)
	}

	return AdToAdResponse(changedAd), nil
}

func (s *Server) UpdateAd(ctx context.Context, req *UpdateAdRequest) (*AdResponse, error) {
	updatedAd, updateErr := s.App.UpdateAd(ctx, req.Id, req.AuthorId, req.Title, req.Text)
	if updateErr != nil {
		return nil, ErrorToGRPCError(updateErr)
	}

	return AdToAdResponse(updatedAd), nil
}

func (s *Server) GetAdByID(ctx context.Context, req *GetAdByIDRequest) (*AdResponse, error) {
	gotAd, getErr := s.App.GetAdByID(ctx, req.Id)
	if getErr != nil {
		return nil, ErrorToGRPCError(getErr)
	}

	return AdToAdResponse(gotAd), nil
}

func (s *Server) GetAdsList(ctx context.Context, req *FilterRequest) (*ListAdResponse, error) {
	date := req.Date
	if date == nil {
		date = &Date{}
	}
	f := filter.Filter{
		PublishedBy: req.PublishedBy,
		AuthorBy:    req.AuthorBy,
		AuthorID:    req.AuthorId,
		DateBy:      req.DateBy,
		Date: ads.Date{
			Day:   int(date.Day),
			Month: time.Month(date.Month),
			Year:  int(date.Year),
		},
	}
	listAd, listErr := s.App.GetAdsList(ctx, f)
	if listErr != nil {
		return nil, ErrorToGRPCError(listErr)
	}

	return AdsListToListAdsResponse(listAd), nil
}

func (s *Server) SearchAds(ctx context.Context, req *SearchAdsRequest) (*ListAdResponse, error) {
	listAd, searchErr := s.App.SearchAds(ctx, req.Pattern)
	if searchErr != nil {
		return nil, ErrorToGRPCError(searchErr)
	}

	return AdsListToListAdsResponse(listAd), nil
}

func (s *Server) DeleteAd(ctx context.Context, req *DeleteAdRequest) (*empty.Empty, error) {
	return &empty.Empty{}, ErrorToGRPCError(s.App.DeleteAd(ctx, req.AuthorId, req.Id))
}
