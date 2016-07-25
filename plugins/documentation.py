# Documentation extension for SpaceDock Backend
# Generates webdocs using python docstrings
# Copyright (c) 2016 The SpaceDock Team
# Licensed under the terms of the MIT License
#
# Name: SpaceDock.extras.documentation

from flask import redirect, url_for, jsonify
from SpaceDock.routing import add_wrapper, route

# Add docs wrapper
methods = {}

# Document functions
def add_documentation(f):
    if not f.__doc__ == None:
        methods[f.api_path] = f.__doc__.strip()
    else:
        methods[f.api_path] = 'No documentation available'
    return f

add_wrapper(add_documentation, True)

# Docs page
@route('/documentation')
def documentation_page():
        return jsonify(methods)