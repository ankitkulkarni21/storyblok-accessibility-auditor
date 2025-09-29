import { createApp } from 'vue'
import App from './App.vue'
import { apiPlugin, StoryblokVue } from "@storyblok/vue";
import './styles.css';

const app = createApp(App);

app.use(StoryblokVue, {
  accessToken: import.meta.env.VITE_STORYBLOK_TOKEN,
  use:[apiPlugin],
  apiOptions: {
    version: "draft",
  },
  bridge: true,
});

app.mount('#app');