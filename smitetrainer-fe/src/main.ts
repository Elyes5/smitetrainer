import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import './themes/theme.css';
import PrimeVue from 'primevue/config';
import store from './store';
createApp(App).use(store).use(router).use(PrimeVue).mount('#app');
