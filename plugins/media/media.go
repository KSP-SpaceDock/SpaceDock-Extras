/*
 Data transformation extension for SpaceDock Backend

 SpaceDock-Extras is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package media

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/middleware"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/objects"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/routes"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
    "github.com/kennygrant/sanitize"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "path/filepath"
    "strings"
    "strconv"
)

func init() {
    routes.Register(routes.POST, "/api/mods/:gameshort/:modid/update-media",
        middleware.NeedsPermission("mods-edit", true, "gameshort", "modid"),
        updateModMedia,
    )
    routes.Register(routes.POST, "/api/users/:userid/update-media",
        middleware.NeedsPermission("user-edit", true, "userid"),
        updateModMedia,
    )
}

/*
 Path: /api/mods/:gameshort/:modid/update-media
 Method: POST
 Description: Uploads a mod background image. Required fields: type, filename, offsetX, offsetY
 Abilities: mods-edit
 */
func updateModMedia(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    mediatype := cast.ToString(utils.GetJSON(ctx, "type"))
    filename := cast.ToString(utils.GetJSON(ctx, "filename"))
    offsetX := cast.ToInt(utils.GetJSON(ctx, "offsetX"))
    offsetY := cast.ToString(utils.GetJSON(ctx, "offsetY"))

    // Check the mediatype
    if mediatype != "background" {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The media type is invalid.").Code(2133))
        return
    }
    ext := filepath.Ext(filepath.Base(filename))
    if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("This file type is not acceptable.").Code(3035))
        return
    }

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

    // Get the file and save it to disk
    filename = sanitize.BaseName(mod.Name) + "_" + mediatype + ext
    base_path := filepath.Join(sanitize.BaseName(mod.User.Username) + "_" + strconv.Itoa(int(mod.User.ID)), sanitize.BaseName(mod.Name))

    // Create a token
    token := objects.NewToken()
    token.SetValue("isUploading", true)
    token.SetValue("path", filepath.Join(base_path, filename))
    app.Database.Save(token)

    // Set data
    mod.SetValue(mediatype, iris.Map{"offsetX": offsetX, "offsetY": offsetY, "path": strings.Replace(filepath.Join(base_path, filename), "\\", "/", -1)})
    app.Database.Save(mod)
    utils.ClearModCache(gameshort, modid)

    // Answer
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": iris.Map{"token": token.Token, "mod": utils.ToMap(mod)}})
}

/*
 Path: /api/users/:userid/update-media
 Method: POST
 Description: Uploads a user background image. Required fields: type, filename, offsetX, offsetY
 Abilities: user-edit
 */
func updateUserMedia(ctx *iris.Context) {
    // Get params
    userid := cast.ToUint(ctx.GetString("userid"))
    mediatype := cast.ToString(utils.GetJSON(ctx, "type"))
    filename := cast.ToString(utils.GetJSON(ctx, "filename"))
    offsetX := cast.ToInt(utils.GetJSON(ctx, "offsetX"))
    offsetY := cast.ToString(utils.GetJSON(ctx, "offsetY"))

    // Check the mediatype
    if mediatype != "background" {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The media type is invalid.").Code(2133))
        return
    }
    ext := filepath.Ext(filepath.Base(filename))
    if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("This file type is not acceptable.").Code(3035))
        return
    }

    // Get the user
    user := &objects.User{}
    app.Database.Where("id = ?", userid).First(user)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The userid is invalid").Code(2130))
        return
    }

    // Get the file and save it to disk
    filename = sanitize.BaseName(user.Username) + "_" + mediatype + ext
    base_path := sanitize.BaseName(user.Username) + "_" + strconv.Itoa(int(user.ID))

    // Create a token
    token := objects.NewToken()
    token.SetValue("isUploading", true)
    token.SetValue("path", filepath.Join(base_path, filename))
    app.Database.Save(token)

    // Set data
    user.SetValue(mediatype, iris.Map{"offsetX": offsetX, "offsetY": offsetY, "path": strings.Replace(filepath.Join(base_path, filename), "\\", "/", -1)})
    app.Database.Save(user)
    utils.ClearUserCache(userid)

    // Answer
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": iris.Map{"token": token.Token, "user": utils.ToMap(user)}})
}
