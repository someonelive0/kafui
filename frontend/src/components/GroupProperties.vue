<template>
  <v-card class="pa-1 ma-1" variant="tonal">
    <v-card-item>
      <div>
        <div class="text-h6 mb-1">
          Consumer Group: {{ name }}
        </div>
        <div class="mb-1">
          GroupID: {{ groupdesc.GroupID }}
        </div>
        <div class="mb-1">
          GroupState: {{ groupdesc.GroupState }}
        </div>
        <div class="mb-1">
          Error: {{ groupdesc.Error }}
        </div>
        <div class="mb-1">
          Members: {{ groupdesc.Members }}
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
import { defineProps, onMounted, ref } from "vue";

const { name } = defineProps(['name']) // 可以简写 解构
let dialog = ref(false);
let snackbar = ref(false);
let snacktext = '';
let groupdesc = ref({
  "Error": null,
  "GroupID": "",
  "GroupState": "",
  "Members": null
});


onMounted(() => {
  refresh();
})

const refresh = () => {

  window.go.backend.KafkaTool.GetGroupDesc(name).then((desc: Uint8Array) => {
    // console.log('Kafkatool.GetGroupDesc ', desc);
    // desc return []byte, need use base64.decode(window.atob) to decode it
    // console.log('Kafkatool.GetGroupDesc ', window.atob(desc));
    if (desc != null) {
      // let descobj = JSON.parse(window.atob(desc));
      groupdesc.value = JSON.parse(window.atob(desc));
    }
  })
  .catch((err: string) => {
    console.error('Kafkatool.GetGroupDesc ', err);
  });
}

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
