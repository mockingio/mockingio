import {defineStore} from 'pinia'
import axios from 'axios'

export interface Mock {
    id: string
    name: string
    routes: Route[]
}

export interface Route {
    request: string
    responses: Response[]
}

export interface Response {
    status: number
    body: string
}

export const useMockStore = defineStore({
    id: "mocks",
    state: () => ({
        mocks: [] as Mock[],
        error: null as any,
    }),
    getters: {
        getMocks: (state) => state.mocks
    },
    actions: {
        async fetchMocks() {
            try {
                const {data} = await axios.get<Mock[]>(getUrl("/mocks"))
                this.mocks = data
            } catch (error) {
                this.error = {
                    error,
                }
            }
        }
    }
})

const getUrl = (path: string): string => {
    return import.meta.env.VITE_API_ENDPOINT + path
}