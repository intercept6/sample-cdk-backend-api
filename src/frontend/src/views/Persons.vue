<template>
  <div class="ion-page">
    <ion-header>
      <ion-toolbar>
        <ion-title>Persons</ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
      <ion-list
          v-for="person in persons"
          v-bind:key="person.Id">
        <ion-item>
          <ion-label>{{ person.FirstName }} {{ person.LastName }}</ion-label>
          <ion-button @click="deleteUser(person.Id)" full>delete</ion-button>
        </ion-item>
      </ion-list>
    </ion-content>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: "about",
  data() { 
    return {
      persons: null
    }
  },
  watch:{
    $route: "reload"
  },
  async created() {
      await this.reload()
  },
  methods: {
    async reload() {
      if(this.$route.fullPath == '/persons'){
        try {
          const response = await axios.get("https://an85ia44y9.execute-api.ap-northeast-1.amazonaws.com/prod/persons")
          console.log(response.data)
          this.persons = response.data
        }catch(e){
          console.log(e)
        }
      }
    },
    async deleteUser(id) {
        console.log(id)
        try {
          await axios.delete("https://an85ia44y9.execute-api.ap-northeast-1.amazonaws.com/prod/persons/"+id)
        }catch(e){
          console.log(e)
        }
        await this.reload()
    }
  }
};
</script>
