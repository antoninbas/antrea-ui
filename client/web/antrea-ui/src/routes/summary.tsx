import React, { useState, useEffect} from 'react';
import { useNavigate } from "react-router-dom";
import { CdsCard } from '@cds/react/card';
import { CdsDivider } from '@cds/react/divider';
import { AgentInfo, ControllerInfo, K8sRef, agentInfoAPI, controllerInfoAPI } from '../api/info';
import { useAccessToken } from '../api/token';
import { APIError } from '../api/common';
import { useAPIError} from '../components/errors';

type Property = string

const controllerProperties: Property[] = ["Name", "Version", "Pod Name", "Node Name", "Connected Agents"]
const agentProperties: Property[] = ["Name", "Version", "Pod Name", "Node Name", "Local Pods", "OVS Version"]

function refToString(ref: K8sRef): string {
    if (ref.namespace) return ref.namespace + '/' + ref.name
    return ref.name
}

function controllerPropertyValues(controller: ControllerInfo): string[] {
    return [
        controller.metadata.name,
        controller.version,
        refToString(controller.podRef),
        refToString(controller.nodeRef),
        (controller.connectedAgentNum??0).toString(),
    ]
}

function agentPropertyValues(agent: AgentInfo): string[] {
    return [
        agent.metadata.name,
        agent.version,
        refToString(agent.podRef),
        refToString(agent.nodeRef),
        (agent.localPodNum??0).toString(),
        agent.ovsInfo.version,
    ]
}

function ComponentSummary<T>(props: {title: string, data: T[], propertyNames: Property[], getProperties: (x: T) => string[]}) {
    const propertyNames = props.propertyNames
    const data = props.data

    return (
        <CdsCard>
            <div cds-layout="vertical gap:md">
                <div cds-text="section" cds-layout="p-y:sm">
                    {props.title}
                </div>
                <CdsDivider cds-card-remove-margin></CdsDivider>
                <table cds-table="border:all" cds-text="center">
                    <thead>
                        <tr>
                            {
                                propertyNames.map(name => (
                                    <th key={name}>{name}</th>
                                ))
                            }
                        </tr>
                    </thead>
                    <tbody>
                        {
                            data.map((x: T, idx: number) => {
                                const values = props.getProperties(x)
                                return (
                                    <tr key={idx}>
                                        {
                                            values.map((v: string, idx: number) => (
                                                <td key={idx}>{v}</td>
                                            ))
                                        }
                                    </tr>
                                )
                            })
                        }
                    </tbody>
                </table>
            </div>
        </CdsCard>
    );
}

export default function Summary() {
    const [controllerInfo, setControllerInfo] = useState<ControllerInfo>();
    const [agentInfos, setAgentInfos] = useState<AgentInfo[]>([]);
    const [accessToken, _] = useAccessToken();
    const { addError } = useAPIError();

    async function getControllerInfo() {
        try {
            const controllerInfo = await controllerInfoAPI.fetch(accessToken)
            setControllerInfo(controllerInfo)
        } catch(e) {
            if (e instanceof APIError) addError(e)
            else throw e
        }
    }

    async function getAgentInfos() {
        try {
            const agentInfos = await agentInfoAPI.fetchAll(accessToken)
            setAgentInfos(agentInfos)
        } catch(e) {
            if (e instanceof APIError) addError(e)
            else throw e
        }
    }

    useEffect(() => {
        getControllerInfo()
        getAgentInfos()
    }, [])

    if (!controllerInfo || !agentInfos) {
        return (
            <p>Loading</p>
        );
    }

    return (
        <main>
            <div cds-layout="vertical gap:lg">
                <p cds-text="title">Summary</p>
                <ComponentSummary title="Controller" data={new Array(controllerInfo)} propertyNames={controllerProperties} getProperties={controllerPropertyValues} />
                <ComponentSummary title="Agents" data={agentInfos} propertyNames={agentProperties} getProperties={agentPropertyValues} />
            </div>
        </main>
    );
}
