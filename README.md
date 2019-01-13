# git-stats

Terminal git stats utility in go.

# Compile

```shell
$ make build
$ make install
```

# Install

After building you can add the folder containing ```git-stats``` to the PATH variable. And now when you invoke ```git stats``` git will look at for the executable named ```git-stats``` and runs it. Thus giving an illusion that ```stats``` is a subcommand of the git.

## Flags

| Flags     | Description                         |
| --summary | Default flag, shows the git summary |
| --ui      | show git commit dots                |
| ...       | ...                                 |
