import {createRouter, createWebHistory} from 'vue-router'
import MockView from '../views/MockView.vue'
import SettingView from '../views/SettingView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: MockView
        },
        {
            path: '/setting',
            name: 'setting',
            component: SettingView
        }
    ]
})

export default router
