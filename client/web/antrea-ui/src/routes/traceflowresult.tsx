import React, { useState, useEffect, useRef} from 'react';
import { useLocation } from "react-router-dom";
import { TraceflowSpec, TraceflowStatus, TraceflowNodeResult } from '../api/traceflow';
import * as d3 from 'd3';
import { graphviz } from "d3-graphviz";

class Node {
    name: string
    attrs: Map<string, string>

    constructor(name: string) {
        this.name = name
        this.attrs = new Map<string, string>()
    }

    setAttr(name: string, value:string) {
        this.attrs.set(name, value)
    }

    asDot(): string {
        const attrs = new Array<string>()
        this.attrs.forEach((v, k) => attrs.push(`${k}=${v}`))
        return `${this.name} [${attrs.join(',')}]`
    }
}

class Edge {
    startNode: string
    endNode: string
    attrs: Map<string, string>

    constructor(startNode: string, endNode: string) {
        this.startNode = startNode
        this.endNode = endNode
        this.attrs = new Map<string, string>()
    }

    asDot(): string {
        const attrs = new Array<string>()
        this.attrs.forEach((v, k) => attrs.push(`${k}=${v}`))
        return `${this.startNode} -> ${this.endNode} [${attrs.join(',')}]`
    }
}

class DotStringBuilder {
    lines: string[]
    indent: number

    constructor() {
        this.lines = new Array<string>()
        this.indent = 0
    }

    addIndent() {
        this.indent += 1
    }

    removeIndent() {
        this.indent -= 1
    }

    pushLine(line: string) {
        const indent = '\t'.repeat(this.indent)
        this.lines.push(indent + line)
    }

    emit(): string {
        return this.lines.join('\n')
    }
}

class Graph {
    graphType: string
    name: string
    nodes: Node[]
    edges: Edge[]
    subgraphs: Graph[]
    attrs: Map<string, string>

    constructor(graphType: string, name: string) {
        this.graphType = graphType
        this.name = name
        this.nodes = new Array<Node>()
        this.edges = new Array<Edge>()
        this.subgraphs = new Array<Graph>()
        this.attrs = new Map<string, string>()
    }

    addNode(node: Node) {
        this.nodes.push(node)
    }

    addEdge(edge: Edge) {
        this.edges.push(edge)
    }

    addSubgraph(graph: Subgraph) {
        this.subgraphs.push(graph)
    }

    setAttr(name: string, value:string) {
        this.attrs.set(name, value)
    }

    asDotBuilder(builder: DotStringBuilder) {
        builder.pushLine(`${this.graphType} ${this.name} {`)
        builder.addIndent()
        this.attrs.forEach((v, k) => builder.pushLine(`${k}=${v}`))
        this.subgraphs.forEach(g => {
            g.asDotBuilder(builder)
            builder.pushLine('')
        })
        this.nodes.forEach(n => builder.pushLine(n.asDot()))
        this.edges.forEach(e => builder.pushLine(e.asDot()))
        builder.removeIndent()
        builder.pushLine('}')
    }

    asDot(): string {
        const builder = new DotStringBuilder()
        this.asDotBuilder(builder)
        return builder.emit()
    }
}

class Digraph extends Graph {
    constructor(name: string) {
        super('digraph', name)
    }
}

class Subgraph extends Graph {
    constructor(name: string) {
        super('subgraph', name)
    }
}

function isSender(nodeResult: TraceflowNodeResult): boolean {
    if (nodeResult.observations.length === 0) {
        return false
    }
    const firstObservation = nodeResult.observations[0]
    if (firstObservation.component !== "SpoofGuard" || firstObservation.action !== "Forwarded") {
        return false
    }
    return true
}

function isReceiver(nodeResult: TraceflowNodeResult): boolean {
    if (nodeResult.observations.length === 0) {
        return false
    }
    const firstObservation = nodeResult.observations[0]
    if (firstObservation.component !== "Forwarding" || firstObservation.action !== "Received") {
        return false
    }
    return true
}

function TraceflowGraph(props: {spec: TraceflowSpec, status: TraceflowStatus}) {
    const tfSpec = props.spec
    const tfStatus = props.status
    const divRef = useRef<HTMLDivElement>(null)

    const darkRed = `"#B20000"`
    const mistyRose = `"#EDD5D5"`
    const fireBrick = `"#B22222"`
    const ghostWhite = `"#F8F8FF"`
    const gainsboro = `"#DCDCDC"`
    const lightGrey = `"#C8C8C8"`
    const silver = `"#C0C0C0"`
    const grey = `"#808080"`
    const dimGrey = `"#696969"`

    useEffect(() => {
        renderGraph(buildGraph())
    });

    function buildGraph(): Digraph {
        console.log(tfStatus)
        const graph = new Digraph('tf')

        if (!tfStatus) return graph

        const senderNodeResult = tfStatus.results.find(isSender)
        const receiverNodeResult = tfStatus.results.find(isReceiver)

        if (!senderNodeResult) return graph

        const srcNode = buildEndpointNode('source', getSourceLabel())
        const [srcCluster, srcLastNode] = buildSubgraph('cluster_source', srcNode, senderNodeResult, false)
        graph.addSubgraph(srcCluster)

        if (!receiverNodeResult) {
            const dstNode = buildEndpointNode('dest', getDestinationLabel())
            srcCluster.addNode(dstNode)
            srcCluster.addEdge(new Edge(srcLastNode.name, dstNode.name))
        }
        // other case (2 nodes) not implemented yet

        return graph
    }

    function getSourceLabel(): string {
        const source = tfSpec.source
        if (source.ip) return source.ip
        return source.namespace + '/' + source.pod
    }

    function getDestinationLabel(): string {
        const dest = tfSpec.destination
        if (dest.ip) return dest.ip
        if (dest.service) return dest.namespace + '/' + dest.service
        return dest.namespace + '/' + dest.pod
    }

    function buildEndpointNode(name: string, label: string): Node {
        const n = new Node(name)
        n.setAttr('style', `"filled,bold"`)
        n.setAttr('label', `"${label}"`)
        n.setAttr('color', grey)
        n.setAttr('fillcolor', lightGrey)
        return n
    }

    function buildSubgraph(name: string, endpointNode: Node, nodeResult: TraceflowNodeResult, needsReversing: boolean): [Subgraph, Node] {
        const graph = new Subgraph('cluster_source')
        graph.setAttr('style', `"filled,bold"`)
        graph.setAttr('bgcolor', ghostWhite)
        graph.setAttr('label', `"${nodeResult.node}"`)
        const nodes = new Array<Node>()
        nodes.push(endpointNode)
        nodeResult.observations.forEach(obs => {
            const n = new Node(obs.component)
            const label = [obs.component, obs.componentInfo, obs.action]
            n.setAttr('shape', `"box"`)
            n.setAttr('style', `"rounded,filled,solid"`)
            n.setAttr('label', `"${label.join('\n')}"`)
            n.setAttr('color', grey)
            n.setAttr('fillcolor', lightGrey)
            nodes.push(n)
        })
        if (needsReversing) nodes.reverse()
        nodes.forEach(n => graph.addNode(n))
        for (let i = 0; i < nodes.length - 1; i++) {
            graph.addEdge(new Edge(nodes[i].name, nodes[i+1].name))
        }
        return [graph, nodes[nodes.length-1]]
    }

    function renderGraph(graph: Digraph) {
        console.log(graph.asDot())
        graphviz(divRef.current).renderDot(graph.asDot());
    }

    return (
        <div ref={divRef}></div>
    );
}

export interface TraceflowResultState {
    spec: TraceflowSpec
    status: TraceflowStatus
}

export default function TraceflowResult() {
    const { state } = useLocation()

    if (!state) return <p>Missing Traceflow Result</p>

    return (
        <div cds-layout="vertical gap:lg">
            <p cds-text="title">Result</p>
            <TraceflowGraph spec={state.spec} status={state.status} />
        </div>
    );
}
