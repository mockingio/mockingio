import {mapState} from "pinia";
import {useMockStore} from "@/stores";

export default {
    ...mapState(useMockStore, ["mocks"]),
}