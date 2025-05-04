
import typing
from taggedunion import *
from epl.chapter5 import continuations

def trampoline(bounce_or_value):
    if type(value) is function:
        return trampoline(value)
    else:
        return value

