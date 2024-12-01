# Info

Helper for maintaining my packages in the AUR.

# Dependencies

- [Go](https://go.dev/) >= 1.23.3
- [Git](https://git-scm.com/)

# Build

Use `make` to build the application.

```shell
make build
```

Then install it in the directory you want.

```shell
make DESTDIR=~/.local/bin install
```

# Configuration

Create the configuration file.

```shell
cp .env.example ~/.config/aur-pkg-helper.env
```

Here's the explanation for each configuration key.

|       Key        | Description                                                                                                |
| :--------------: | :--------------------------------------------------------------------------------------------------------- |
| `AUR_ROOT_PATH`  | The path that holds all the AUR package git repositories. Non-git directories will be ignored.             |
| `GIT_HOOKS_PATH` | The path where will be stored the Git hooks that will be setup on each git repository of the AUR packages. |
| `GIT_USER_NAME`  | Git user name to be used on each git repository.                                                           |
| `GIT_USER_EMAIL` | Git user email to be setup on each git repository.                                                         |

> **Note:** The paths defined in `AUR_ROOT_PATH` and `GIT_HOOKS_PATH` must exist and must have enought permissions to read and write files to.
