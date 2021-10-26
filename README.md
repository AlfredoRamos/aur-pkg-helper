### Info

A C++ app to manage git hooks for my AUR packages.

### Dependencies

- `git` >= 5.0.0
- `cmake` >= 3.18.0
- C++ compilet with `c++20` support


### Submodules

```shell
git submodule update --init --remote
```

### Build

```shell
cmake -S . -B build \
	-DCMAKE_INSTALL_PREFIX=/usr \
	-DCMAKE_BUILD_TYPE=Release
cmake --build build --clean-first
```

### Install

```shell
cmake --install build --strip
```
