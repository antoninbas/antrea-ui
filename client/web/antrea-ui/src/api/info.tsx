import { APIError, handleErrorResponse } from './common'
import config from '../config';
const { apiServer, apiUri } = config;

export interface K8sRef {
    namespace?: string
    name: string
}

interface NetworkPolicyControllerInfo {
    networkPolicyNum: number
    addressGroupNum: number
    appliedToGroupNum: number
}

interface ControllerCondition {
    type: string
    status: string
    lastHeartbeatTime: string
    reason: string
    message: string
}

export interface ControllerInfo {
    metadata: {
        name: string
    }
    version: string
    podRef: K8sRef
    nodeRef: K8sRef
    serviceRef: K8sRef
    networkPolicyControllerInfo: NetworkPolicyControllerInfo
    connectedAgentNum: number
    controllerConditions: ControllerCondition[]
    apiPort: number
}

interface OVSInfo {
    version: string
    bridgeName: string
    flowTable: Map<string,number>
}

interface AgentCondition {
    type: string
    status: string
    lastHeartbeatTime: string
    reason: string
    message: string
}

export interface AgentInfo {
    metadata: {
        name: string
    }
    version: string
    podRef: K8sRef
    nodeRef: K8sRef
    nodeSubnets: string[]
    ovsInfo: OVSInfo
    networkPolicyControllerInfo: NetworkPolicyControllerInfo
    localPodNum: number
    agentConditions: AgentCondition[]
    apiPort: number
}

export const controllerInfoAPI = {
    fetch: async (token: string): Promise<ControllerInfo> => {
        try {
            const response = await fetch(`${apiUri}/info/controller`, {
                method: "GET",
                mode: "cors",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
            });

            return handleErrorResponse(response).then((data) => data as ControllerInfo);
        } catch (err) {
            console.error("Unable to fetch Controller Info");
            throw err;
        }
    },
}

export const agentInfoAPI = {
    fetchAll: async (token: string): Promise<AgentInfo[]> => {
        try {
            const response = await fetch(`${apiUri}/info/agents`, {
                method: "GET",
                mode: "cors",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
            });

            return handleErrorResponse(response).then((data) => data as AgentInfo[]);
        } catch (err) {
            console.error("Unable to fetch Agent Infos");
            throw err;
        }
    },

    fetch: async (name: string, token: string): Promise<AgentInfo> => {
        try {
            const response = await fetch(`${apiUri}/info/agents/${name}`, {
                method: "GET",
                mode: "cors",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
            });

            return handleErrorResponse(response).then((data) => data as AgentInfo);
        } catch (err) {
            console.error("Unable to fetch Agent Info");
            throw err;
        }
    },
}
