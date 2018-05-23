package mysql

import (
	"log"

	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// Scores retrieves all scores from the DB.
func (db DB) Scores() ([]clong.Score, error) {
	var scrs []clong.Score
	rows, err := db.Query("SELECT id, playerID, playerName, finalScore, color FROM scores")
	if err != nil {
		return []clong.Score{}, errors.Wrap(err, "error getting scores from DB")
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error closing DB rows"))
		}
	}()
	for rows.Next() {
		var s clong.Score
		err = rows.Scan(&s.ID, &s.Player.ID, &s.Player.Name, &s.FinalScore, &s.Color)
		if err != nil {
			return []clong.Score{}, errors.Wrap(err, "error scanning DB rows")
		}
		scrs = append(scrs, s)
	}
	err = rows.Err()
	if err != nil {
		return []clong.Score{}, errors.Wrap(err, "error in DB rows")
	}

	return scrs, nil
}
