<template>
  <v-container fluid class="pa-1 ma-1">
    <v-row align="center" justify="left">
      <v-col cols="auto">
        <v-card
          append-icon="mdi-check"
          class="mx-auto"
          max-width="344" min-width="300"
          prepend-icon="mdi-server"
          subtitle="Number of brokers"
          title="Brokers"
          hover
        >
          <v-card-text>Total: {{ num_brokers }}</v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row align="center" justify="left">
      <v-col cols="auto">
        <v-card
          class="mx-auto"
          max-width="344" min-width="300"
          subtitle="Number of topics"
          title="Topics"
          hover
        >
          <template v-slot:prepend>
            <v-icon color="primary" icon="mdi-format-list-bulleted"></v-icon>
          </template>
          <template v-slot:append>
            <v-icon color="success" icon="mdi-check"></v-icon>
          </template>
          <v-card-text>Total: {{ num_topics }}</v-card-text>
          <v-card-actions>
            <v-btn color="teal-accent-4" variant="text" @click="gotoTopics">
              Learn More
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <v-row align="center" justify="left">
      <v-col cols="auto">
        <v-card
          class="mx-auto"
          max-width="344" min-width="300"
          subtitle="Number of consumer groups"
          title="Consumer Groups"
          hover
        >
          <template v-slot:prepend>
            <v-icon color="primary" icon="mdi-account-multiple"></v-icon>
          </template>
          <template v-slot:append>
            <v-icon color="success" icon="mdi-check"></v-icon>
          </template>
          <v-card-text>Total: {{ num_groups }}</v-card-text>
          <v-card-actions>
            <v-btn color="teal-accent-4" variant="text" @click="gotoGroups">
              Learn More
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter(); 
const { query, params } = useRoute();
// console.log('{ query, params } = useRoute() ', query, params);

const num_brokers = ref(query.num_brokers);
const num_topics = ref(query.num_topics);
const num_groups = ref(query.num_groups);
// const brokers = query.brokers;
const topics = query.topics;
const groups = query.groups;


// const gotoBrokers = () => {
//   router.push({
//     name: 'Brokers',
//     query: { brokers: brokers }
//   });
// }

const gotoTopics = () => {
  // router.push('Topics');
  router.push({
    name: 'Topics',
    query: { topics: topics }
  });
}

const gotoGroups = () => {
  router.push({
    name: 'Groups',
    query: { groups: groups }
  });
}

</script>
