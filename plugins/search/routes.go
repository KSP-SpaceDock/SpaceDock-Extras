/*
 Search extension for SpaceDock Backend

 SpaceDock-Extras is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package search

import (
    "SpaceDock"
    "SpaceDock/middleware"
    "SpaceDock/objects"
    "SpaceDock/routes"
    "SpaceDock/utils"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "sort"
)

func init() {
    routes.Register(routes.GET, "/api/browse/:gameshort", browse_mod)
    routes.Register(routes.GET, "/api/browse/:gameshort/:mode", browse_mod_mode)
}

/*
 Path: /api/browse/:gameshort
 Method: GET
 Description: Sorts all mods of the game by popularity.
 */
func browse_mod(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    count := ctx.URLParam("count")
    site := ctx.URLParam("site")

    // Get Data
    data, status := grabMods(ctx, gameshort, site, count)
    utils.WriteJSON(ctx, status, data)
}

/*
 Path: /api/browse/:gameshort/:mode
 Method: GET
 Description: Sorts all mods of the game by popularity.
 */
func browse_mod_mode(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    mode := ctx.GetString("mode")
    count := ctx.URLParam("count")
    site := ctx.URLParam("site")

    // Get Data
    data, status := grabMods(ctx, gameshort, site, count)

    // Errorcheck
    if _,ok := data[mode]; !ok {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Invalid mode.").Code(3900))
        return
    }
    utils.WriteJSON(ctx, status, data[mode])
}

type ByUpdated struct {
    values *[]objects.Mod
}

func (s ByUpdated) Len() int {
    return len(*(s.values))
}
func (s ByUpdated) Swap(i, j int) {
    (*(s.values))[i], (*(s.values))[j] = (*(s.values))[j], (*(s.values))[i]
}
func (s ByUpdated) Less(i, j int) bool {
    return (*(s.values))[i].UpdatedAt.Before((*(s.values))[j].UpdatedAt)
}

func forEach(data []objects.Mod, f func(interface{}) map[string]interface{}) []interface{} {
    result := []interface{}{}
    for _,e := range data {
        result = append(result, f(e))
    }
    return result
}

func forEachFeatured(data []objects.Featured, f func(interface{}) map[string]interface{}) []interface{} {
    result := []interface{}{}
    for _,e := range data {
        result = append(result, f(e))
    }
    return result
}

func grabMods(ctx *iris.Context, gameshort string, site_ string, count_ string) (iris.Map, int) {
    game := &objects.Game{}
    SpaceDock.Database.Where("short = ?", gameshort).First(game)
    if game.Short != gameshort {
        return utils.Error("The game does not exist.").Code(2125), iris.StatusNotFound
    }

    // Params
    site := 0
    if s,err := cast.ToIntE(site_); err == nil {
        site = s - 1
    }
    count := 6
    if c,err := cast.ToIntE(count_); err == nil {
        count = c
    }

    // Magical numbers
    magic := count * site
    if magic > 0 {
        magic = magic - 1
    }

    // Get the mods
    featured := []objects.Featured{}
    SpaceDock.Database.Joins("JOIN mods ON mod.id = featured.mod_id").
        Where("mods.game_id = ?", game.ID).
        Order("featured.created_at DESC").
        Find(&featured)
    featured = featured[magic:count]
    top,_ := searchMods(game, "", float64(site + 1), count)
    top = top[:count]
    new := []objects.Mod{}
    SpaceDock.Database.Where("published = ?", true).Where("game_id = ?", game.ID).Order("created_at DESC").Find(&new)
    new = new[magic:count]
    updated := []objects.Mod{}
    SpaceDock.Database.Where("published = ?", true).
        Where("game_id = ?", game.ID).
        Where("created_at != updated_at").
        Order("updated_at DESC").
        Find(&new)
    updated = updated[magic:count]
    current_user := middleware.CurrentUser(ctx)
    yours := current_user.Following
    if current_user != nil {
        sort.Sort(sort.Reverse(ByUpdated{values:&yours}))
        yours = yours[magic:count]
    }
    data := iris.Map {
        "featured": forEachFeatured(featured, utils.ToMap),
        "top": forEach(top, utils.ToMap),
        "new": forEach(new, utils.ToMap),
        "updated": forEach(updated, utils.ToMap),
        "yours": forEach(yours, utils.ToMap),
    }
    return data, iris.StatusOK
}