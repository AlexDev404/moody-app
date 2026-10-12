package playlist

import "testing"

func TestParse(t *testing.T) {
	songs, err := parse("```json\n{\"name\":\"generate_playlist\",\"parameters\":{\"mood\":\"calm\",\"playlist\":[" +
		"{\"title\":\"Weightless\",\"artist\":\"Marconi Union\",\"mood_tags\":[\"calm\"],\"reason\":\"slow" + string(rune(0x00A0)) + "and soft\"}," +
		"{\"title\":\"\",\"artist\":\"Nobody\"}]}}\n```")
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 1 || songs[0].Title != "Weightless" || songs[0].Reason != "slow and soft" {
		t.Errorf("songs = %+v", songs)
	}

	for _, bad := range []string{`{"status": 1}`, `not json`, `{"parameters":{"playlist":[]}}`} {
		if _, err := parse(bad); err == nil {
			t.Errorf("parse(%q) succeeded", bad)
		}
	}
}
