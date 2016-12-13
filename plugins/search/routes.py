# Search extension for SpaceDock Backend
# Copyright (c) 2016 The SpaceDock Team
# Licensed under the terms of the MIT License
#
# Name: SpaceDock.extras.search.routes
# Depends: search/algorithm.py

from flask import request
from flask_login import current_user
from sqlalchemy import desc
from SpaceDock.formatting import feature_info, mod_info
from SpaceDock.objects import Featured, Game, Mod
from SpaceDock.routing import route
from SpaceDock.extras.search.algorithm import search_mods, search_users, typeahead_mods


@route("/api/browse/<gameshort>")
def browse_mod(gameshort):
    """
    Sorts all mods of the game by popularity.
    """
    # get params
    site = request.args.get('site')
    count = request.args.get('count')
    
    return grab_mods(gameshort, site, count)
    
@route("/api/browse/<gameshort>/<mode>")
def browse_mod_mode(gameshort, mode):
    """
    Sorts all mods of the game by popularity.
    """
    # get params
    site = request.args.get('site')
    count = request.args.get('count')

    # Grab Data
    data = grab_mods(gameshort, site, count)
    
    # Errorcheck
    if not mode in data['data']:
        return {'error': True, 'reasons': ['Invalid mode.']}, 400
        
    # Assemble the return
    data = data['data'][mode]
    return {'error': False, 'count': count, 'data': data}
    
def grab_mods(gameshort, site, count):
    """
    Grabs mods from the databased, sorted by featured, popularity, new and updated
    """
    # Validate gameshort
    if not Game.query.filter(Game.short == gameshort).first():
        return {'error': True, 'reasons': ['The gameshort is invalid.']}, 400
        
    # Params
    if not site:
        site = 0
    else:
        site = int(site) - 1
    if not count:
        count = 6
    else:
        count = int(count)
        
    # Magical numbers
    magic = count * site
    if magic > 0: magic = magic - 1
    
    # Get the game
    ga = Game.query.filter(Game.short == gameshort).first()
    featured = Featured.query.outerjoin(Mod).filter(Mod.game_id == ga.id).order_by(desc(Featured.created)).all()[magic:count]
    top = search_mods(ga,"", site + 1, count)[:count][0]
    new = Mod.query.filter(Mod.published, Mod.game_id == ga.id).order_by(desc(Mod.created)).all()[magic:count]
    updated = Mod.query.filter(Mod.published, Mod.game_id == ga.id, Mod.updated != Mod.created).order_by(desc(Mod.updated)).all()[magic:count]
    yours = []
    if current_user:
        yours = sorted(current_user.following, key=lambda m: m.updated, reverse=True)[magic:count]
    data = {
        'featured': [feature_info(f) for f in featured],
        'top': [mod_info(m) for m in top],
        'new': [mod_info(m) for m in new],
        'updated': [mod_info(m) for m in updated],
        'yours': [mod_info(m) for m in yours]
    }
    return {'error': False, 'count': count, 'data': data}