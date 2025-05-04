# epl

Python implementation of exercises in the Elements of Programming Languages (Friedman and Wand, 3rd Ed)

## Requirements

1. pyenv (https://github.com/pyenv/pyenv-virtualenv)

2. Python 3.0

```
pyenv install 3.7.0
```

3. Create env:

```
pyenv virtualenv 3.7.0 epl
pyenv activate epl
```

## Setup

```
pip install -r requirements
```

## Testing

```
cd tests
pytest -s
```
