<template>
    <v-container class="grey lighten-5">
      <v-row class="mb-6" no-gutters >
        <v-col key="1" cols="4" >
          <v-card class="pa-2" outlined tile >
            <v-text-field
                v-model="search"
                clear-icon="mdi-close-circle-outline"
                label="Search Company Directory"
                clearable
                dark
                flat
                hide-details
                solo-inverted
            ></v-text-field>
            <v-treeview
                v-model="tree"
                :items="items"
                :search="search"
                :opened="initiallyOpen"
                item-key="name"
                activatable
                open-on-click
                density="compact"
            >
                <template v-slot:prepend="{ item, open }">
                <v-icon v-if="!item.file">
                    {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
                </v-icon>
                <v-icon v-else>
                    {{ files[item.file] }}
                </v-icon>
                </template>
            </v-treeview>
          </v-card>
        </v-col>

        <v-col key="2" cols="8" >
          <v-card class="pa-2" outlined tile >
            xx
          </v-card>
        </v-col>
      </v-row>
      
    </v-container>
  </template>

<script setup lang="ts">
import { ref, reactive, defineProps } from "vue"
import { VTreeview } from 'vuetify/labs/VTreeview'


let search = ref('');
const initiallyOpen = ['public'];
const files = {
    html: 'mdi-language-html5',
    js: 'mdi-nodejs',
    json: 'mdi-code-json',
    md: 'mdi-language-markdown',
    pdf: 'mdi-file-pdf-box',
    png: 'mdi-file-image',
    txt: 'mdi-file-document-outline',
    xls: 'mdi-file-excel',
};
const tree = ref([]);
const items = [
    {
        title: '.git',
    },
    {
        title: 'node_modules',
    },
    {
        title: 'public',
        children: [
        {
            title: 'static',
            children: [{
            title: 'logo.png',
            file: 'png',
            }],
        },
        {
            title: 'favicon.ico',
            file: 'png',
        },
        {
            title: 'index.html',
            file: 'html',
        },
        ],
    },
    {
        title: '.gitignore',
        file: 'txt',
    },
    {
        title: 'babel.config.js',
        file: 'js',
    },
    {
        title: 'package.json',
        file: 'json',
    },
    {
        title: 'README.md',
        file: 'md',
    },
    {
        title: 'vue.config.js',
        file: 'js',
    },
    {
        title: 'yarn.lock',
        file: 'txt',
    },
];


</script>
