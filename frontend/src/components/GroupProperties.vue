<template>
  <v-card class="pa-1 ma-1" variant="tonal">
    <v-card-item>
      <div>
        <div class="text-overline mb-1">
          ID: {{ name }}
        </div>
        <div class="text-h6 mb-1">
          Host: {{ name }}
        </div>
      </div>
    </v-card-item>

    <v-card-actions>
      <v-dialog v-model="dialog" max-width="500" >
        <template v-slot:activator="{ props: activatorProps }">
          <v-btn color="red-darken-4" variant="outlined" v-bind="activatorProps">Delete This Group</v-btn>
        </template>

        <v-card
          prepend-icon="mdi-map-marker"
          text="This consumer group will be deleted, and lost all information of the group."
          :title="'Really delete group ' + name + ' ?'" 
        >
          <template v-slot:actions>
            <v-spacer></v-spacer>
            <v-btn @click="dialog = false">Cancel</v-btn>&nbsp;
            <v-btn color="red-darken-4" @click="deleteGroup">Delete</v-btn>
          </template>
        </v-card>

      </v-dialog>
    </v-card-actions>
  </v-card>

  <v-snackbar v-model="snackbar" timeout=2000 color="deep-purple-accent-4" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="pink" variant="text" @click="snackbar = false" >Close</v-btn>
    </template>
  </v-snackbar>

</template>

<script setup lang="ts">
import { ref, defineProps } from "vue"

const { name } = defineProps(['name']) // 可以简写 解构
let dialog = ref(false);
let snackbar = ref(false);
let snacktext = '';

const deleteGroup = () => {
  console.log('deleteGroup ', name);
  window.go.backend.KafkaTool.DeleteGroup(name).then(() => {
    snacktext = 'delete group ' + name + ' success!';
    snackbar.value = true;
    dialog.value = false;
  })
  .catch((err: string) => {
    // console.error('Kafkatool.WriteMsg ', err);
    snacktext = 'delete group ' + name + ' failed: ' + err;
    snackbar.value = true;
  });
  
}

</script>
