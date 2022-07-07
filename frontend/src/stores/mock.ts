import {defineStore} from 'pinia'
import axios from 'axios'

export enum Status {
    STOPPED = "stopped",
    RUNNING = "running"
}

export interface MockData {
    id: string
    name: string
    description: string
    routes: Route[]
}

export interface MockState {
    url?: string
    status: Status
}

export interface Mock {
    data: MockData
    state: MockState
}

export interface Mocks {
    [id: string]: Mock
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
        activeId: null as string | null,
        mocks: {} as Mocks,
        error: null as any,
    }),
    getters: {
        getMockByID(state) {
            return (id: string): Mock | undefined => state.mocks[id]
        },
        activeMock(state): Mock | undefined {
            if (!!state.activeId) {
                return this.getMockByID(state.activeId)
            }

            const ids = Object.keys(state.mocks)
            if (ids.length === 0) {
                return undefined
            }

            this.activeId = ids[0]
            return this.getMockByID(this.activeId)
        }
    },
    actions: {
        setActiveMock(id: string) {
            this.activeId = id
        },
        async fetchMocks() {
            try {
                const {data} = await axios.get<MockData[]>(getUrl("/mocks"))
                this.mocks = data.reduce((acc, mock) => {
                    acc[mock.id] = {
                        data: mock,
                        state: {
                            status: Status.STOPPED
                        }
                    }
                    return acc
                }, {} as Mocks)
            } catch (error) {
                this.error = {
                    error,
                }
            }
        },
    }
})

const getUrl = (path: string): string => {
    return import.meta.env.VITE_API_ENDPOINT + path
}