# Wails + Vue 3 Typescript

[English](README.md) | 简体中文

Dependencies
Wails has a number of common dependencies that are required before installation:

Go 1.18+
NPM (Node 15+)

```shell
go install github.com/wailsapp/wails/v2/cmd/wails@latest

wails init -n "ProjectName" -t https://github.com/airvip/wails-vite-vue-ts.git
```

## 关于

这是一个 Wails 模板项目,使用 Vue 3 和 TypeScript,使用 Vite 为 asset bundling。
它提供了最基本的东西,可以根据这篇 README 的指导进行扩展。

如果你想要包含更多特性,请查看我的 [feature packed Vite + Vue3 TypeScript template](https://github.com/airvip/wails-vite-vue-the-works)


## 实施开发

执行开发模式, 运行 `wails dev` 在项目目录. 在其他终端，进入 `frontend`
目录并运行 `npm run dev`. 通过浏览器打开 http://localhost:34115 连接到你的应用。

说明：类型检查已经被关闭。如果你想要进行类型检查,请使用 `npm run type-check`

即，开一个终端，进入 `frontend`目录并运行 `npm run dev`. 然后再开一个终端并运行 `wails dev`，
表示 `npm run dev`实时跟踪 `frontend/src`目录下文件变化，然后编译更新到 `frontend/dist`目录下，
而 `wails dev`实时跟踪 `frontend/dist`目录下文件变化，不过需要在浏览器手工F5刷新页面，这样就可以做到前端的开发模式。
避免每次更新前端源代码都执行 `wails dev`


## 扩展特性

这个模板不包括路由、vuex 或 sass。
增加这些功能，只需简单地按照下面的指导进行操作。
请注意所有指令应该在 `frontend` 目录中运行。

### Sass

安装:
```shell
$ npm install --save-dev sass
```

使用:

你可以添加 Sass 到单个文件组件样式如下:
```html
<style lang="scss">
  /* scss styling */
</style>
```

### ESLint + Prettier

安装:
```shell
$ npm install --save-dev eslint prettier eslint-plugin-vue eslint-config-prettier @vue/eslint-config-typescript @typescript-eslint/parser @typescript-eslint/eslint-plugin
$ touch .eslintrc && touch .prettierrc
```

使用: `eslintrc`
```json
{
  "extends": [
    "plugin:vue/vue3-essential",
    "eslint:recommended",
    "prettier",
    "@vue/typescript/recommended"
  ],
    "rules": {
    // override/add rules settings here, such as:
    // "vue/no-unused-vars": "error"
  }
}
```

使用: `.prettierrc`
```json
{
  "semi": false,
  "tabWidth": 2,
  "useTabs": false,
  "printWidth": 120,
  "endOfLine": "auto",
  "singleQuote": true,
  "trailingComma": "all",
  "bracketSpacing": true,
  "arrowParens": "always"
}
```

### Vuex

安装:
```shell
$ npm install --save vuex@next
$ touch src/store.ts
```

使用: `src/store.ts`
```ts
import { InjectionKey } from 'vue'
import { createStore, Store, useStore as baseUseStore } from 'vuex'

// define your typings for the store state
export interface State {
  count: number
}

// define injection key
export const key: InjectionKey<Store<State>> = Symbol()

export const store = createStore<State>({
  state() {
    return {
      count: 0
    }
  },
  mutations: {
    increment(state) {
      state.count++
    }
  }
})

export function useStore() {
  return baseUseStore(key)
}
```

使用: `src/main.ts`
```ts
import { createApp } from 'vue'
import App from './App.vue'
import { store, key } from './store'

createApp(App).use(store, key).mount('#app')
```

使用: `src/components/Home.vue`
```ts
import { useStore } from '../store'
const store = useStore()
const increment = () => store.commit('increment')
```

### Vue Router

安装:
```shell
$ npm install --save vue-router@4
$ touch src/router.ts
```

使用: `src/router.ts`
```ts
import { createRouter, createWebHashHistory } from 'vue-router'
import Home from './components/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
```

使用: `src/main.ts`
```ts
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

createApp(App).use(router).mount('#app')
```

使用: `src/App.vue`
```html
<template>
    <router-link to="/">Home</router-link>
    <router-view />
</template>
```

## 编译 

编译这个项目在 debug 模式，使用 `wails build`. 在生产模式下，使用 `wails build -production`。
生产某个平台原生包，需要使用 `-package` 选项。

## 已知问题

- 当你在开发模式下运行时,改变了前端页面，浏览器会非常快速的自动刷新页面，导致问题。刷新页面能解决。
- 类型检查已经被关闭是因为 Wails 依赖于前端构建之前进行构建,生成绑定。
- 如果你发现了其他问题,请创建一个 issue。