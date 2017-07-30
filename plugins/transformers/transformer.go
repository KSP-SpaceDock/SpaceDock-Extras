/*
 Data transformation extension for SpaceDock Backend

 SpaceDock-Extras is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package transformers

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/objects"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
)

func init() {
    utils.RegisterDataTransformer(Transform)
}

func Transform(data interface{}, m map[string]interface{}) {
    if mod, ok := data.(*objects.Mod); ok {
        m["follower_count"] = len(mod.Followers)
        m["author"] = mod.User.Username
        m["game_short"] = mod.Game.Short
    }
    if mod, ok := data.(objects.Mod); ok {
        m["follower_count"] = len(mod.Followers)
        m["author"] = mod.User.Username
        m["game_short"] = mod.Game.Short
    }
    if _, ok := data.(*objects.Featured); ok {
        mod := &objects.Mod{}
        app.Database.Where("id = ?", m["mod_id"]).First(mod)
        (m["mod"].(map[string]interface{}))["follower_count"] = len(mod.Followers)
        (m["mod"].(map[string]interface{}))["author"] = mod.User.Username
        (m["mod"].(map[string]interface{}))["game_short"] = mod.Game.Short
    }
    if _, ok := data.(objects.Featured); ok {
        mod := &objects.Mod{}
        app.Database.Where("id = ?", m["mod_id"]).First(mod)
        (m["mod"].(map[string]interface{}))["follower_count"] = len(mod.Followers)
        (m["mod"].(map[string]interface{}))["author"] = mod.User.Username
        (m["mod"].(map[string]interface{}))["game_short"] = mod.Game.Short
    }
}
