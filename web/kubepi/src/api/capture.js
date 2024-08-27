import {get, post, put} from "@/plugins/request"


const baseUrl = "/api/v1/clusters"

export function getTaskAll() {
    return get(`${baseUrl}/captureall`)
}

export function setTask(data) {
    return post(`${baseUrl}/capturetask`,data)
}


export function getConfig() {
    return get(`${baseUrl}/captureconfig`)
}

export function setConfig(data) {
    return post(`${baseUrl}/captureconfig`, data)
}

