package api

import (
	"context"
	"math"
	"time"

	"google.golang.org/appengine/log"

	"google.golang.org/appengine/datastore"
)

const DAKind = "DeviationAdjustment"

var laddersToSave []*Ladder
var ladderKeysToSave []*datastore.Key

var playersToSave []*Player
var playerKeysToSave []*datastore.Key

// DeviationAdjustment represents a calculation of a player's deviation
// in line with step 6 in http://www.glicko.net/glicko/glicko2.pdf
type DeviationAdjustment struct {
	Date       time.Time
	Player     Player
	Adjustment float64
}

func (da *DeviationAdjustment) Save(ctx context.Context) (*datastore.Key, error) {
	key := datastore.NewIncompleteKey(ctx, DAKind, nil)
	return datastore.Put(ctx, key, da)
}

func NewDeviationAdjustment(p *Player) DeviationAdjustment {
	return DeviationAdjustment{
		Date:       time.Now(),
		Player:     *p,
		Adjustment: math.Sqrt(math.Pow(p.RatingDeviation, 2) + math.Pow(p.RatingVolatility, 2)),
	}
}

func (l *Ladder) MinActivitityDate() time.Time {
	return time.Now().Add(-l.InactivityPeriod)
}

func Active(ctx context.Context, l *Ladder, p *Player) (bool, error) {
	last, err := p.LastActive(ctx)

	if err != nil {
		log.Errorf(ctx, "[inactivity] error getting player last active time: %s (player %v)", err.Error(), p)
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
		for j, lp := range ladder.Players {
			p, err := PlayerFromLadderPlayer(ctx, lp)

			if err != nil {
				log.Errorf(ctx, "[inactivity] could not create player from ladder player %v: %s", lp, err.Error())
			}

			active, err := Active(ctx, ladder, p)

			if err != nil {
				return err
			}

			if !active {
				// skip deviation update for now, needs testing
				// updateDeviation(ctx, p, ladder)

				if lp.Active {
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

func updateDeviation(ctx context.Context, p *Player, l *Ladder) error {
	lastAdjustment := make([]*DeviationAdjustment, 0)

	if _, err := datastore.NewQuery(LadderKind).
		Filter("Player.FirebaseID = ", p.FirebaseID).
		Order("-Date").
		Limit(1).
		GetAll(ctx, &lastAdjustment); err != nil {
		log.Errorf(ctx, "[inactivity] could not query DeviationAdjustment: %s", err.Error())
		return err
	}

	if len(lastAdjustment) == 0 || lastAdjustment[0].Date.Add(l.InactivityPeriod).Before(time.Now()) {
		da := NewDeviationAdjustment(p)

		p.RatingDeviation = p.RatingDeviation + da.Adjustment

		playersToSave = append(playersToSave, p)
		playerKeysToSave = append(playerKeysToSave, p.DatastoreKey(ctx))

		if _, err := da.Save(ctx); err != nil {
			log.Errorf(ctx, "[inactivity] could not save DeviationAdjustment: %s", err.Error())
			return err
		}
	}

	return nil
}
