<template>
    <v-card flat
      class="mx-auto"
      prepend-icon="mdi-copyright"
      width="400"
    >

      <template v-slot:title>
        KAFUI
      </template>
      <template v-slot:subtitle>
        Version: 0.0.5 &nbsp; 2025-08-26
      </template>

      <v-card-text>
        Kafui is a Kafka GUI client.<br/>
        This software is under 
        
        <B @click="openBrower('https://mit-license.org')">
          <v-icon icon="mdi-open-in-new" size="small"></v-icon>&nbsp;MIT license
        </B>.
        <v-img :src="myLogo" width="330" ></v-img>
      </v-card-text>

      <v-card-actions>
        <v-btn color="secondary" text="GITHUB" @click="openBrower('https://github.com/someonelive0/kafui')"></v-btn>
        <v-spacer></v-spacer>
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
