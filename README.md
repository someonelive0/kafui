# KAFUI

Kafui is a Kafka GUI client. This software is under MIT license .

Kafui use Wails 2 + Vue 3 Typescript.
And use wails-vite-vue-ts from https://github.com/airvip/wails-vite-vue-ts.git

Backend use kafka-go from https://github.com/segmentio/kafka-go

## License

This software is under MIT license, to visit file [LICENSE](LICENSE) and https://mit-license.org/ for details.


## Build Enviroment 

English | [简体中文](README.zh-CN.md) show how to develop this project

Depend golang-1.23, nodejs-20 or later.

```shell
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails doctor

# build frontend
cd frontend
npm run build

# build app
wails build
wails build --tags exp_gowebview2loader
```

kafui.exe will found in build/bin


## Live Development

```shell
# run frontend
cd frontend
npm run dev

# run backend
wails dev
```

Then visit http://localhost:34115 to debug this project.


To run in live development mode, run `wails dev` in the project directory. In another terminal, go into the `frontend`
directory and run `npm run dev`. Navigate to http://localhost:34115
in your browser to connect to your application.

Note: Typechecking is disabled. If you want to do type checking, use `npm run type-check`


## Config

Run kafui.exe, app read local kafui.toml, if it is not existed copy from kafui.toml.tpl.

Open 'Setting' icon on toolbar top set kafka brokers config items. Then test connection.

Flow is from kafui.toml.tpl, show how kafka config items in toml format file.

```
[kafka]
    name = "localhost"
    brokers = [ "127.0.0.1:9092" ]
    # sasl mechanism should be empty or "SASL_PLAINTEXT",
    # if mechanism is "SASL_PLAINTEXT", then set user and password
    sasl_mechanism = ""
    user = ""
    password = ""
```
