<template>
  <v-card flat>
    <v-card-title class="d-flex align-center pe-2">
      <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
      {{ name }} Offsets list &nbsp; {{ offsets.length }}
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Search"
        prepend-inner-icon="mdi-magnify"
        variant="outlined"
        hide-details
        single-line
        density="compact"
      ><v-tooltip activator="parent" location="bottom">Match keyword</v-tooltip>
      </v-text-field>&nbsp;
      <v-btn icon="mdi-refresh" size="small" @click="refresh"></v-btn>&nbsp;
    </v-card-title>


    <v-data-table
      :headers="headers"
      :items="offsets"
      :search="search"
      :loading="loading"
      :items-per-page="-1"
      density="compact"
      item-value="topic"
      hover
    >
      <template v-slot:item.actions="{ item }">
        <div class="d-flex ga-2 justify-end">
          <v-icon color="medium-emphasis" icon="mdi-pencil" size="small" @click="edit(item)"></v-icon>
        </div>
      </template>

    </v-data-table>


    <v-dialog v-model="dialog" max-width="500">
      <v-card
        :subtitle="`${isEditing ? 'Set' : 'Create'} new commited offset of consumer group`"
        :title="`${isEditing ? 'Edit' : 'Add'} new commited offset of consumer group`"
      >
        <template v-slot:text>
          <v-row>
            <v-col cols="6">
              <v-text-field v-model="formModel.topic" label="Topic" disabled></v-text-field>
            </v-col>

            <v-col cols="6" md="6">
              <v-text-field v-model="formModel.partition" label="Partition" disabled></v-text-field>
            </v-col>

            <v-col cols="6" md="6">
              <v-text-field v-model="formModel.first_offset" label="First Offset" disabled></v-text-field>
            </v-col>

            <v-col cols="6" md="6">
              <v-text-field v-model="formModel.last_offset" label="Last Offset" disabled></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field v-model="formModel.committed_offset" label="New Offset" variant="outlined"></v-text-field>
            </v-col>
            <v-col cols="12">
              <small class="text-caption text-medium-emphasis">*New Offset must &gt; 'First Offset' and &lt; 'Last Offset'</small>
            </v-col>
          </v-row>
        </template>

        <v-divider></v-divider>

        <v-card-actions class="bg-surface-light">
          <v-btn text="Cancel" variant="plain" @click="dialog = false"></v-btn>

          <v-spacer></v-spacer>

          <v-btn text="Save" @click="save"></v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>

  <v-snackbar v-model="snackbar" timeout=2000 color="deep-purple-accent-4" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="pink" variant="text" @click="snackbar = false" >Close</v-btn>
    </template>
  </v-snackbar>

</template>


<script setup lang="ts">
import { defineProps, onMounted, reactive, ref, shallowRef, toRef } from 'vue';
import { backend } from '../wailsjs/go/models';


const { name } = defineProps(['name']) // 可以简写 解构

const headers = [
      { title: 'Topic', align: 'start', sortable: true, key: 'topic' },
      { title: 'Partition', align: 'end', key: 'partition' },
      { title: 'First Offset', align: 'end', key: 'first_offset' },
      { title: 'Last Offset', align: 'end', key: 'last_offset' },
      { title: 'Committed Offset', align: 'end', key: 'committed_offset' },
      { title: 'Actions', align: 'end', key: 'actions' },
    ];

let offsets: Array<backend.GroupOffset> = reactive([]);
let loading = ref(true);
let search = ref('');
const sortBy = [{ key: 'config_name', order: 'asc' }];
let snackbar = ref(false);
let snacktext = '';

// for edit group offset
const formModel = ref([])
const dialog = shallowRef(false)
const isEditing = toRef(() => !!formModel.value.topic)


onMounted(() => {
  refresh();
})

const refresh = () => {
  loading.value = true; // why not work?

  window.go.backend.KafkaTool.GetGroupOffset(name).then((items: Array<backend.GroupOffset>) => {
    // console.log('Kafkatool.GetGroupOffset ', items);
    if (items != null) {
      offsets = items;
    }
    loading.value = false;
  })
  .catch((err: string) => {
    console.error('Kafkatool.GetGroupOffset ', err);
    loading.value = false;
  });
}

const edit = (item) => {
  // const found = books.value.find(book => book.id === id)
  console.log("edit item: " + item.topic);
  formModel.value = item;
  dialog.value = true
}

const save = () => {
  console.log("save item: " + formModel.value.topic);
  formModel.value.committed_offset = formModel.value.committed_offset.trim();
  if (formModel.value.committed_offset.length == 0 || isNaN(Number(formModel.value.committed_offset)) ) {
    snacktext = 'Error: New Offset must be numeric!';
    snackbar.value = true;
    return;
  }
  if (parseInt(formModel.value.committed_offset) < parseInt(formModel.value.first_offset)
    || parseInt(formModel.value.committed_offset) > parseInt(formModel.value.last_offset) ) {
    snacktext = 'Error: New Offset ' + formModel.value.committed_offset + ' is less than [First Offset] or bigger than [Last Offset]!';
    snackbar.value = true;
    return;
  }

  window.go.backend.KafkaTool.SetGroupOffset(name, formModel.value.topic, 
    formModel.value.partition, 
    parseInt(formModel.value.committed_offset) ).then(() => {
    snacktext = 'set group offset ' + name + ' success!';
    snackbar.value = true;

    dialog.value = false
    refresh();
  })
  .catch((err: string) => {
    // console.error('Kafkatool.SetGroupOffset ', err);
    snacktext = 'set group offset ' + name + ' failed: ' + err;
    snackbar.value = true;
  });

}

</script>
