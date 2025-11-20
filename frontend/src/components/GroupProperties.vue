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
          text="This consumer group will be deleted, and will lost all information of the group."
          :title="'Really delete Group [' + name + '] ?'" 
        >
          <template v-slot:prepend>
            <v-icon color="red" icon="mdi-alert"></v-icon>
          </template>
          <template v-slot:actions>
            <v-spacer></v-spacer>
            <v-btn border @click="dialog = false">Cancel</v-btn>&nbsp;
            <v-btn border color="red-darken-4" @click="deleteGroup">Delete</v-btn>
          </template>
        </v-card>

      </v-dialog>
    </v-card-actions>
  </v-card>

  <v-snackbar v-model="snackbar" timeout=4000 :color="snackcolor" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="grey" variant="text" @click="snackbar = false">Close</v-btn>
    </template>
  </v-snackbar>

</template>


<script setup lang="ts">
import { defineProps, onMounted, ref } from "vue";
import { GetGroupDesc, DeleteGroup } from "../wailsjs/go/backend/KafkaTool";


const { name } = defineProps(['name']) // 可以简写 解构
let dialog = ref(false);
let snackbar = ref(false);
let snacktext = '';
let snackcolor = 'deep-purple-darken-4';

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

  GetGroupDesc(name).then((desc: string) => {
    // console.log('Kafkatool.GetGroupDesc ', desc);
    groupdesc.value = JSON.parse(desc);
  }).catch((err: string) => {
    console.error('Kafkatool.GetGroupDesc failed: ', err);
    groupdesc.value.GroupState = "BUSY";
  });
}

const deleteGroup = () => {
  console.log('deleteGroup ', name);
  DeleteGroup(name).then(() => {
    showSnackBar('delete group ' + name + ' success!', true);
    dialog.value = false;
  }).catch((err: string) => {
    showSnackBar('delete group ' + name + ' failed: ' + err, false);
  });
  
}

const showSnackBar = (text: string, success: boolean) => {
    snackbar.value = false;
    snacktext = text;
    snackcolor = success ? 'deep-purple-darken-4' : 'deep-orange-darken-3';
    snackbar.value = true;
}
</script>
