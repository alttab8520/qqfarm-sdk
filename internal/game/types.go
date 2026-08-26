package game

import (
	"context"
	"errors"
)

var ErrNotLogin = errors.New("未登录")

type User struct {
	GID      int64  `json:"gid"`
	Name     string `json:"name"`
	Level    int64  `json:"level"`
	Exp      int64  `json:"exp"`
	Gold     int64  `json:"gold"`
	OpenID   string `json:"open_id"`
	Avatar   string `json:"avatar_url"`
}

type Land struct {
	ID         int64  `json:"id"`
	Unlocked   bool   `json:"unlocked"`
	Level      int64  `json:"level"`
	PlantID    int64  `json:"plant_id,omitempty"`
	PlantName  string `json:"plant_name,omitempty"`
	FruitID    int64  `json:"fruit_id,omitempty"`
	DryNum     int64  `json:"dry_num,omitempty"`
	HasWeed    bool   `json:"has_weed"`
	HasInsect  bool   `json:"has_insect"`
	Stealable  bool   `json:"stealable"`
	LeftFruit  int64  `json:"left_fruit,omitempty"`
}

type Item struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
}

type Friend struct {
	GID    int64  `json:"gid"`
	Name   string `json:"name"`
	Level  int64  `json:"level"`
	OpenID string `json:"open_id"`
	Avatar string `json:"avatar_url"`
}

type LoginIn struct {
	Code   string `json:"code"`
	OpenID string `json:"open_id"`
}

type HarvestIn struct {
	LandIDs []int64 `json:"land_ids"`
	HostGID int64   `json:"host_gid"`
	IsAll   bool    `json:"is_all"`
}

type PlantIn struct {
	SeedID  int64   `json:"seed_id"`
	LandIDs []int64 `json:"land_ids"`
}

type HelpIn struct {
	GID     int64   `json:"gid"`
	LandIDs []int64 `json:"land_ids"`
}

type LandOpIn struct {
	LandIDs []int64 `json:"land_ids"`
	HostGID int64   `json:"host_gid"`
}

type FertilizeIn struct {
	LandIDs      []int64 `json:"land_ids"`
	FertilizerID int64   `json:"fertilizer_id"`
}

type Status struct {
	LoggedIn bool   `json:"logged_in"`
	User     User   `json:"user,omitempty"`
	ACE      ACE    `json:"ace"`
}

type ACE struct {
	Uploads   int    `json:"uploads"`
	Reports   int    `json:"status_reports"`
	Failures  int    `json:"failures"`
	LastError string `json:"last_error,omitempty"`
}

// Session is the live game connection used by HTTP handlers.
type Session interface {
	Login(ctx context.Context, in LoginIn) (User, error)
	Info() (User, error)
	Status() (Status, error)
	Refresh(ctx context.Context) ([]Land, error)
	Harvest(ctx context.Context, in HarvestIn) ([]Item, error)
	Plant(ctx context.Context, in PlantIn) error
	Water(ctx context.Context, in LandOpIn) error
	Weed(ctx context.Context, in LandOpIn) error
	Bug(ctx context.Context, in LandOpIn) error
	Fertilize(ctx context.Context, in FertilizeIn) error
	Friends(ctx context.Context) ([]Friend, error)
	Help(ctx context.Context, in HelpIn) error
	Close() error
}

type Factory func() Session
