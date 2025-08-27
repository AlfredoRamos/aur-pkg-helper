# Info

Helper for maintaining my packages in the AUR.

# Dependencies

- [Go](https://go.dev/) >= 1.25.0
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

Create the configuration file `~/.config/aur-pkg-helper.toml`

```toml
[aur]
root_path = 

[git]
hooks_path = 
user_name = 
user_email = 
```

Here's the explanation for each configuration key.

|       Key        | Description                                                                                                |
| :--------------: | :--------------------------------------------------------------------------------------------------------- |
| `aur.root_path`  | The path that holds all the AUR package git repositories. Non-git directories will be ignored.             |
| `git.hooks_path` | The path where will be stored the Git hooks that will be setup on each git repository of the AUR packages. |
| `git.user_name`  | Git user name to be used on each git repository.                                                           |
| `git.user_email` | Git user email to be setup on each git repository.                                                         |

> **Note:** The paths defined in `aur.root_path` and `git.hooks_path` must exist and must have enought permissions to read and write files to.
