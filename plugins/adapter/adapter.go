/*
 Route Adapters for SpaceDock Backend

 SpaceDock-Extras is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package adapter

import (
    "SpaceDock"
    "SpaceDock/routes"
    "SpaceDock/utils"
    "SpaceDock/objects"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
)

func init() {
    routes.Register(routes.GET, "/api/adapter/mods/:modid", mods_adapter)
}

/*
 Path: /api/adapter/mods/:modid
 Method: GET
 Description: Returns information for one mod
 */
func mods_adapter(ctx *iris.Context) {
    // Get params
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mods gameshort
    mod := &objects.Mod{}
    SpaceDock.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    ctx.Redirect("/api/mods/" + mod.Game.Short + "/" + cast.ToString(modid), iris.StatusPermanentRedirect)
}