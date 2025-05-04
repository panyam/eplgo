
from ipdb import set_trace

class Settings(object):
    default_settings = None

    def __init__(self):
        self.stack = [{}]

    @classmethod
    def defaults(cls):
        if cls.default_settings is None:
            cls.default_settings = Settings()
        return cls.default_settings

    def __enter__(self):
        return self

    def __exit__(self, type, value, traceback):
        self.pop()
        return False

    def push(self, **kwargs):
        self.stack.append({})
        for k,v in kwargs.items():
            self.stack[-1][k] = v
        return self

    def pop(self):
        self.stack.pop()

    def get(self, key, no_throw = False):
        for frame in reversed(self.stack):
            if key in frame:
                return frame[key]
        if no_throw: return None
        raise AttributeError(key)

def push(**kwargs):
    return Settings.defaults().push(**kwargs)

def get(key, no_throw = False):
    defaults = Settings.defaults()
    return defaults.get(key, no_throw = no_throw)
