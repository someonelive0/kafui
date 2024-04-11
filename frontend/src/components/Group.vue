<template>
  <v-container fluid class="pa-1 ma-1">
    <v-row no-gutters>
      <v-col>
        <v-breadcrumbs :items="breadcrumbs" class="bg-grey-lighten-5 pa-1 ma-1">
          <template v-slot:prepend>
            <v-icon icon="mdi-account-file-text-outline" size="small"></v-icon>
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
        <v-tab value="Offset">Offset</v-tab>
      </v-tabs>
      
      <v-window v-model="selectedTab">
        <v-window-item key="Properties" value="Properties">
          <v-container fluid>
            <GroupProperties :name="group"/>
          </v-container>
        </v-window-item>

        <v-window-item key="Offset" value="Offset">
          <v-container fluid>
            <GroupOffset :name="group"/>
          </v-container>
        </v-window-item>

      </v-window>
    </v-card>

  </v-container>
</template>

<script setup lang="ts">
import { ref, isRef } from 'vue';
import { useRoute } from 'vue-router';
import GroupProperties from './GroupProperties.vue'
import GroupOffset from './GroupOffset.vue'

const { query, params } = useRoute();
// console.log('{ query, params } = useRoute() ', query, params);
const group = ref(query.group);
var breadcrumbs = ref([
  { title: 'Consumer Group', disabled: false, },
  { title: group, disabled: false, }
]);
let selectedTab = ref("Properties"); // 默认选中 Properties 页

</script>
