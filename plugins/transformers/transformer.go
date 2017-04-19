/*
 Data transformation extension for SpaceDock Backend

 SpaceDock-Extras is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package transformers

import (
    "SpaceDock"
    "SpaceDock/objects"
    "SpaceDock/utils"
)

func init() {
    utils.RegisterDataTransformer(Transform)
}

func Transform(data interface{}, m map[string]interface{}) {
    if mod, ok := data.(*objects.Mod); ok {
        m["follower_count"] = len(mod.Followers)
        m["author"] = mod.User.Username
    }
    if mod, ok := data.(objects.Mod); ok {
        m["follower_count"] = len(mod.Followers)
        m["author"] = mod.User.Username
    }
    if _, ok := data.(*objects.Featured); ok {
        mod := &objects.Mod{}
        SpaceDock.Database.Where("id = ?", m["mod_id"]).First(mod)
        (m["mod"].(map[string]interface{}))["follower_count"] = len(mod.Followers)
        (m["mod"].(map[string]interface{}))["author"] = mod.User.Username
    }
    if _, ok := data.(objects.Featured); ok {
        mod := &objects.Mod{}
        SpaceDock.Database.Where("id = ?", m["mod_id"]).First(mod)
        (m["mod"].(map[string]interface{}))["follower_count"] = len(mod.Followers)
        (m["mod"].(map[string]interface{}))["author"] = mod.User.Username
    }
}
