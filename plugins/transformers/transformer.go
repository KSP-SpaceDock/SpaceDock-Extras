/*
 Data transformation extension for SpaceDock Backend

 SpaceDock-Extras is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package transformers

import (
    "SpaceDock/objects"
    "SpaceDock/utils"
)

func init() {
    utils.RegisterDataTransformer(TransformMod)
}

func TransformMod(data interface{}, m map[string]interface{}) {
    if mod, ok := data.(*objects.Mod); ok {
        m["follower_count"] = len(mod.Followers)
        m["author"] = mod.User.Username
    }
}
