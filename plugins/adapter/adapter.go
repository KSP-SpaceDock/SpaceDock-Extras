/*
 Route Adapters for SpaceDock Backend

 SpaceDock-Extras is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package adapter

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/middleware"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/objects"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/routes"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
)

func init() {
    routes.Register(routes.GET, "/api/adapter/mods/:modid", middleware.Recursion(0), mods_adapter)
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
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    app.Database.Model(mod).Related(&(mod.Game), "Game")
    if ctx.URLParam("callback") != "" {
        ctx.Redirect("/api/mods/" + mod.Game.Short + "/" + cast.ToString(modid) + "?callback=" + ctx.URLParam("callback"), iris.StatusPermanentRedirect)
        return
    }
    ctx.Redirect("/api/mods/" + mod.Game.Short + "/" + cast.ToString(modid), iris.StatusPermanentRedirect)
}