<template>
    <v-card flat
      class="mx-auto"
      width="400"
    >
      <template v-slot:prepend>
        <v-avatar color="blue-accent-2">
          <v-icon icon="mdi-copyright"></v-icon>
        </v-avatar>
      </template>
      <template v-slot:title>
        KAFUI
      </template>
      <template v-slot:subtitle>
        Version: 0.0.6 &nbsp; 2025-09-01
      </template>

      <v-card-text>
        Kafui is a Kafka GUI client.<br/>
        This software is under 
        
        <B @click="openBrower('https://mit-license.org')">
          <v-icon icon="mdi-open-in-new" color="blue-darken-3" size="small"></v-icon>&nbsp;MIT license
        </B>.
        <v-img :src="myLogo" width="330" ></v-img>
      </v-card-text>

      <v-card-actions>
        <v-btn density="compact" icon="mdi-github" color="blue-darken-3"
          @click="openBrower('https://github.com/someonelive0/kafui')"></v-btn>
        https://github.com/someonelive0/kafui
      </v-card-actions>

    </v-card>
</template>


<script setup lang="ts">
import myLogo from "../assets/images/appicon.png";
import { Goto } from "../wailsjs/go/main/App";


const openBrower = (url: string) => {
  Goto(url).then(() => { // window.go.main.App.Goto
  })
  .catch((err: string) => {
    console.error('visitPage', err);
  });
}

</script>
