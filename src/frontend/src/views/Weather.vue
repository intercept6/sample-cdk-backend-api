<template>
    <div class="weather">
        <h1>This is an weather page</h1>
    </div>
</template>


<script>
    import axios from 'axios';

    export default {
        name: "weather",
        data() {
            return {
                info: null
            };
        },
        created() {
            const url = 'https://fyn0prjgd6.execute-api.ap-northeast-1.amazonaws.com/prod/things';
            axios.get(url)
                .then(response => (console.log(response.data)))
                .catch(error => (this.info = error));
        }
    }
</script>
