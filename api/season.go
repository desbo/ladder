package api

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func (l *Ladder) SnapshotDatastoreKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, LadderSnapshotKind, "", int64(l.Season), l.DatastoreKey(ctx))
}

func (l *Ladder) SaveSnapshot(ctx context.Context) (*datastore.Key, error) {
	return datastore.Put(ctx, l.SnapshotDatastoreKey(ctx), l)
}

func (l *Ladder) StartNewSeason(ctx context.Context) (*datastore.Key, error) {
	if _, err := l.SaveSnapshot(ctx); err != nil {
		return nil, err
	}

	l.Season = l.Season + 1
	l.SeasonStart = time.Now()
	newPlayers := make(Players, 0)

	for _, p := range l.Players {
		if p.Active {
			p.Wins = 0
			p.Losses = 0
			p.Rating = initialRating
			p.RatingDeviation = 350
			p.RatingVolatility = 0.06
			p.LastRatingChange = 0
			p.WinRate = 0
			newPlayers = append(newPlayers, p)
		}
	}

	l.Players = newPlayers

	return l.Save(ctx)
}
