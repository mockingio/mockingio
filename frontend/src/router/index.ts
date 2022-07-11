import {createRouter, createWebHistory} from 'vue-router'
import MockView from '../views/MockView.vue'
import SettingView from '../views/SettingView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'default',
            component: MockView
        },
        {
            path: '/mocks/:id',
            name: 'mockview',
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
