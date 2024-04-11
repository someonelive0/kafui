<template>
  <v-container fluid class="pa-1 ma-1">
    <v-row no-gutters>
      <v-col>
        <v-breadcrumbs :items="breadcrumbs" class="bg-grey-lighten-5 pa-1 ma-1">
          <template v-slot:prepend>
            <v-icon icon="mdi-book-open-variant-outline" size="small"></v-icon>
          </template>
          <template v-slot:divider>
                <v-icon icon="mdi-chevron-right"></v-icon>
            </template>
        </v-breadcrumbs>
      </v-col>
    </v-row>

    <v-card>
      <v-tabs density="compact"
        v-model="selectedTab"
        align-tabs="start"
        color="deep-purple-accent-4"
        >
        <v-tab value="Properties">Properties</v-tab>
        <v-tab value="Data">Data</v-tab>
        <v-tab value="Config">Config</v-tab>
      </v-tabs>
      
      <v-window v-model="selectedTab">
        <v-window-item key="Properties" value="Properties">
          <v-container fluid>
            <TopicProperties :name="topic" />
          </v-container>
        </v-window-item>

        <v-window-item key="Data" value="Data">
          <v-container fluid>
            <Data :name="topic" />
          </v-container>
        </v-window-item>

        <v-window-item key="Config" value="Config">
          <v-container fluid>
            <Config title="topic" :name="topic" />
          </v-container>
        </v-window-item>

      </v-window>
    </v-card>

  </v-container>
</template>

<script setup lang="ts">
import { ref, isRef } from 'vue';
import { useRoute } from 'vue-router';
import TopicProperties from './TopicProperties.vue'
import Data from './Data.vue'
import Config from './Config.vue'


const { query, params } = useRoute();
// console.log('{ query, params } = useRoute() ', query, params);
const topic = ref(query.topic);
var breadcrumbs = ref([
  { title: 'Topic', disabled: false, },
  { title: topic, disabled: false, }
]);

let selectedTab = ref("Properties"); // 默认选中 Properties 页

</script>
