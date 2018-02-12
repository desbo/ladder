package api

import (
	"context"
	"time"

	"google.golang.org/appengine/log"

	"google.golang.org/appengine/datastore"
)

const DAKind = "DeviationAdjustment"

var laddersToSave []*Ladder
var ladderKeysToSave []*datastore.Key

var playersToSave []*Player
var playerKeysToSave []*datastore.Key

func (l *Ladder) MinActivitityDate() time.Time {
	return time.Now().Add(-l.InactivityPeriod)
}

func Active(ctx context.Context, l *Ladder, u *User) (bool, error) {
	last, err := u.LastActive(ctx)

	if err != nil {
		log.Errorf(ctx, "[inactivity] error getting user last active time: %s (user %v)", err.Error(), u)
		return false, err
	}

	return last.After(l.MinActivitityDate()), nil
}

func CheckInactivity(ctx context.Context) error {
	ladders := make([]*Ladder, 0)

	if _, err := datastore.NewQuery(LadderKind).Filter("Active = ", true).GetAll(ctx, &ladders); err != nil {
		log.Errorf(ctx, "error getting ladders: %s", err.Error())
		return err
	}

	for _, ladder := range ladders {
		for j, p := range ladder.Players {
			u, err := UserFromPlayer(ctx, p)

			if err != nil {
				log.Errorf(ctx, "[inactivity] could not create user from player %v: %s", p, err.Error())
			}

			active, err := Active(ctx, ladder, u)

			if err != nil {
				return err
			}

			if !active {
				if p.Active {
					ladder.Players[j].Active = false
					laddersToSave = append(laddersToSave, ladder)
					ladderKeysToSave = append(ladderKeysToSave, ladder.DatastoreKey(ctx))
				}
			}
		}
	}

	if _, err := datastore.PutMulti(ctx, ladderKeysToSave, laddersToSave); err != nil {
		return err
	}

	if _, err := datastore.PutMulti(ctx, playerKeysToSave, playersToSave); err != nil {
		return err
	}

	return nil
}
