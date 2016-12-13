# Route Adapters for SpaceDock Backend
# Copyright (c) 2016 The SpaceDock Team
# Licensed under the terms of the MIT License
#
# Name: SpaceDock.extras.adapter
from SpaceDock.formatting import mod_info
from SpaceDock.objects import Mod
from SpaceDock.routing import route

@route('/api/adapter/mods/<modid>')
def mods_adapter(modid):
    """
    Returns information for one mod
    """
    if not modid.isdigit() or not Mod.query.filter(Mod.id == int(modid)).first():
        return {'error': True, 'reasons': ['The modid is invalid']}, 400
    # Get the mod
    mod = Mod.query.filter(Mod.id == int(modid)).first()
    if not mod.published and current_user != mod.user:
        return {'error': True, 'reasons': ['The mod is not published.']}, 400
    return {'error': False, 'count': 1, 'data': mod_info(mod)}