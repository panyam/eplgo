
class TypeError(Exception):
    def __init__(self, expected, found):
        Exception.__init__(self, "Expected type: %s, Found: %s" % (expected, type))
        self.expected = expected
        self.found = found

def ensure_type(expected, found):
    if expected != found:
        raise TypeError(expected, found)
