import { APIError, handleErrorResponse } from './common'
import config from '../config';
const { apiServer, apiUri } = config;

interface ObjectMetadata {
    creationTimestamp?: string
    name?: string
}

export interface TraceflowPacket {
    ipHeader: {
        protocol?: number
    }
    transportHeader: {
        icmp?: {
        }
        udp?: {
            srcPort?: number
            dstPort?: number
        }
        tcp?: {
            srcPort?: number
            dstPort?: number
            flags?: number
        }
    }
}

export interface TraceflowSpec {
    source: {
        namespace?: string
        pod?: string
        ip?: string
    }
    destination: {
        namespace?: string
        pod?: string
        service?: string
        ip?: string
    }
    packet?: TraceflowPacket
    timeout?: number
}

export interface TraceflowObservation {
    component: string
    componentInfo: string
    action: string
    pod: string
    dstMAC: string
    networkPolicy: string
    egress: string
    ttl: number
    translatedSrcIP: string
    translatedDstIP: string
    tunnelDstIP: string
    egressIP: string
}

export interface TraceflowNodeResult {
    node: string
    role: string
    timestamp: number
    observations: TraceflowObservation[]
}

export interface TraceflowStatus {
    phase: string
    reason: string
    startTime: string
    results: TraceflowNodeResult[]
}

interface Traceflow {
    apiVersion?: string
    kind?: string
    metadata?: ObjectMetadata
    spec?: TraceflowSpec
    status?: TraceflowStatus
}

export const traceflowAPI = {
    runTraceflow: async (tf: TraceflowSpec, token: string): Promise<TraceflowStatus | undefined> => {
        try {
            let url = `${apiUri}/traceflow`
            let response = await fetch(url, {
                method: "POST",
                mode: "cors",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
                body: JSON.stringify({spec: tf}),
            });

            if (response.status !== 202) {
                throw new APIError(response.status, response.statusText, "Error when creating traceflow request");
            }

            for (let i = 0; i < 10; i++) {
                let location = response.headers.get("Location") ?? "";
                let retryAfter = response.headers.get("Retry-After") ?? "";
                let waitFor = parseInt(retryAfter) * 1000;
                await new Promise(r => setTimeout(r, waitFor));
                url = `${apiServer}${location}`
                response = await fetch(url, {
                    method: "GET",
                    mode: "cors",
                    headers: {
                        "Authorization": `Bearer ${token}`,
                    },
                });
                if (response.status === 200) {
                    return response.json().then((data) => data as Traceflow).then((tf) => tf.status)
                }
                if (response.status !== 202) {
                    throw new APIError(response.status, response.statusText, "Error when checking traceflow request status");
                }
            }
            throw new APIError(0, "", "Timeout when waiting for traceflow request to compleye")
        } catch (err) {
            console.error("Unable to run traceflow");
            throw err;
        }
    }
}
