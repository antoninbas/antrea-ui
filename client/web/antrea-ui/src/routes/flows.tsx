import React, { useState, useEffect} from 'react';
import ReactECharts from "echarts-for-react";
import { CdsCard } from '@cds/react/card';
import { CdsDivider } from '@cds/react/divider';
import flowsRaw from './flows.json'

const flows: Flow[] = flowsRaw;

flows.sort((f1, f2) => f1.startTime - f2.startTime)

interface Flow {
    source: string
    destination: string
    destinationPort: number
    startTime: number
    endTime: number
}

function latestFlows(count: number): Flow[] {
    return flows.slice(-count)
}

function getTimeseriesData(): Array<[string, number]> {
    const data = new Map<string, number>()
    flows.forEach(f => {
        const mins = Math.floor(f.startTime / 60)
        const t = new Date(mins * 60 * 1000).toISOString()
        const c = data.get(t) ?? 0
        data.set(t, c + 1)
    })
    return Array.from(data.entries())
}

// const option = {
//   xAxis: {
//     type: 'category',
//     data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
//   },
//   yAxis: {
//     type: 'value'
//   },
//   series: [
//     {
//       data: [120, 200, 150, 80, 70, 110, 130],
//       type: 'line'
//     }
//   ]
// };

const option = {
  tooltip: {
    trigger: 'axis',
  },
  dataset: {
    source:getTimeseriesData(),
    dimensions: ['timestamp', 'rate'],
  },
  xAxis: {
    type: 'time',
  },
  yAxis: { },
  series: [
    {
        name: 'rate',
        type: 'line',
        encode: {
            x: 'timestamp',
            y: 'rate',
        },
    }
  ]
};

// in order for 100% width to work, we downgrade echarts-for-react to 3.0.1
// we can switch to 3.0.3 once released
// see https://github.com/hustcc/echarts-for-react/pull/464
function FlowsTimeseries() {
    return (
        <CdsCard>
            <div cds-layout="vertical gap:sm">
                <div cds-text="section" cds-layout="p-y:sm">
                    New Flows / Minute
                </div>
                <CdsDivider cds-card-remove-margin></CdsDivider>
                <ReactECharts option={option} style={{ height: '400px', width: '100%' }} />
            </div>
        </CdsCard>
    );
}

function FlowsTable(props: {count: number}) {
    const flows = latestFlows(props.count);

    return (
        <CdsCard>
            <div cds-layout="vertical gap:md">
                <div cds-text="section" cds-layout="p-y:sm">
                    Latest Flows
                </div>
                <CdsDivider cds-card-remove-margin></CdsDivider>
                <table cds-table="border:all" cds-text="center">
                    <thead>
                        <tr>
                            <th>Source</th>
                            <th>Destination</th>
                            <th>Destination Port</th>
                            <th>Start Time</th>
                            <th>End Time</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            flows.map((f, idx) => (
                                <tr key={idx}>
                                    <td>{f.source}</td>
                                    <td>{f.destination}</td>
                                    <td>{f.destinationPort}</td>
                                    <td>{new Date(f.startTime * 1000).toLocaleString()}</td>
                                    <td>{new Date(f.endTime * 1000).toLocaleString()}</td>
                                </tr>
                            ))
                        }
                    </tbody>
                </table>
            </div>
        </CdsCard>
    );
}

export default function FlowVisibility() {
    useEffect(() => {
    }, [])

    return (
        <main>
            <div cds-layout="vertical gap:lg">
                <p cds-text="title">Flow Visibility</p>
                <FlowsTimeseries />
                <FlowsTable count={25} />
            </div>
        </main>
    );
}
