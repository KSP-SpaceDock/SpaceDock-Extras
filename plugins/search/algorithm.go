/*
 Search extension for SpaceDock Backend

 SpaceDock-Extras is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package search

import (
    "SpaceDock"
    "SpaceDock/objects"
    "github.com/spf13/cast"
    "math"
    "sort"
    "strings"
    "time"
)

type ByWeight struct {
    values *[]objects.Mod
    terms []string
}

func (s ByWeight) Len() int {
    return len(*(s.values))
}
func (s ByWeight) Swap(i, j int) {
    (*(s.values))[i], (*(s.values))[j] = (*(s.values))[j], (*(s.values))[i]
}
func (s ByWeight) Less(i, j int) bool {
    return weightResult(&((*(s.values))[i]), s.terms) < weightResult(&((*(s.values))[j]), s.terms)
}

func weightResult(result *objects.Mod, terms []string) float64 {
    /* Factors considered, * indicates important factors:
     * Mods where several search terms match are given a dramatically higher rank*
     * High followers and high downloads get bumped*
     * Mods with a long version history get bumped
     * Mods with lots of screenshots or videos get bumped
     * Mods with a short description get docked
     * Mods lose points the longer they go without updates*
     * Mods get points for supporting the latest KSP version
     * Mods get points for being open source
     * New mods are given a hefty bonus to avoid drowning among established mods
     */
    var score float64
    name_matches := 0
    short_matches := 0
    for _,term := range terms {
        if strings.Count(strings.ToLower(result.Name), term) != 0 {
            name_matches += 1
            score += float64(name_matches * 100)
        }
        if strings.Count(strings.ToLower(result.ShortDescription), term) != 0 {
            short_matches += 1
            score += float64(short_matches * 50)
        }
    }
    score *= 100

    score += float64(len(result.Followers) * 10)
    score += float64(result.DownloadCount)
    score += float64(len(result.Versions)) / float64(5)
    if len(result.Description) < 100 {
        score -= 10
    }
    delta := time.Now().Sub(result.UpdatedAt).Hours() / 24
    if delta > 100 {
        delta = 100 // Don't penalize for oldness past a certain point
    }
    score -= delta / 5
    if err,_ := result.GetValue("source_link"); err == nil {
        score += 10
    }
    if (time.Now().Sub(result.CreatedAt).Hours() / 24) < 30 {
        score += 100
    }
    return score
}

func searchMods(game *objects.Game, text string, page float64, limit int) ([]objects.Mod, float64) {
    terms := strings.Split(text, " ")
    results := []objects.Mod{}
    query := SpaceDock.Database.Joins("JOIN users ON users.id = mods.user_id").
        Joins("JOIN mod_versions ON mod_versions.mod_id = mods.id").
        Joins("JOIN games ON games.id = mods.game_id").
        Joins("JOIN game_versions ON game_versions.id = mod_versions.game_version_id")
    queries := []string{}
    for _,term := range terms {
        if strings.HasPrefix(term, "ver:") {
            queries = append(queries, "game_versions.friendly_version = '" + term[4:] + "'")
        } else if strings.HasPrefix(term, "user:") {
            queries = append(queries, "users.username = '" + term[5:] + "'")
        } else if strings.HasPrefix(term, "game:") {
            queries = append(queries, "mod.game_id = " + cast.ToString(cast.ToInt(term[5:])))
        } else if strings.HasPrefix(term, "downloads:>") {
            queries = append(queries, "mod.download_count > " + cast.ToString(cast.ToInt(term[11:])))
        } else if strings.HasPrefix(term, "downloads:<") {
            queries = append(queries, "mod.download_count < " + cast.ToString(cast.ToInt(term[11:])))
        } else {
            queries = append(queries, "LOWER(mods.name) LIKE '%" + strings.ToLower(term) + "%'")
            queries = append(queries, "LOWER(users.username) LIKE '%" + strings.ToLower(term) + "%'")
            queries = append(queries, "LOWER(mods.short_description) LIKE '%" + strings.ToLower(term) + "%'")
        }
    }
    if game != nil {
        query = query.Where("mods.game_id = ?", game.ID)
    }
    query = query.Where("(" + strings.Join(queries, " OR ") + ")")
    query = query.Where("mods.published = ?", true)
    query.Find(&results)

    total := math.Ceil(float64(len(results)) / float64(limit))
    if page > total {
        page = total
    }
    if page < 1 {
        page = 1
    }
    sort.Sort(sort.Reverse(ByWeight{values:&results,terms:terms}))
    return results[(int(page) - 1) * limit:int(math.Min(float64(int(page) * limit), float64(cap(results))))], total
}

func searchUsers(text string, page float64) []objects.User {
    terms := strings.Split(text, " ")
    results := []objects.User{}
    queries := []string{}
    for _,term := range terms {
        queries = append(queries, "LOWER(users.username) LIKE '%"+strings.ToLower(term)+"%'")
        queries = append(queries, "LOWER(users.description) LIKE '%"+strings.ToLower(term)+"%'")
    }
    query := SpaceDock.Database.Where("(" + strings.Join(queries, " OR ") + ")")
    query = query.Where("users.public = ?", true)
    query.Find(&results)
    return results[(int(page)) * 10:(int(page) * 10) + 10]
}