"use client";

import { useState, useEffect, useRef } from "react";
import Graph from "graphology";
import { circular } from "graphology-layout";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Github, Search } from "lucide-react";

const FAKE_NODES = [
  "Alice", "Bob", "Charlie", "David", "Emma", "Frank",
  "Grace", "Henry", "Ivy", "Jack", "Kate", "Liam",
  "Mia", "Noah", "Olivia", "Paul", "Quinn", "Ruby"
];

const FAKE_EDGES = [
  ["Alice", "Bob"], ["Alice", "Charlie"], ["Bob", "David"],
  ["Charlie", "Emma"], ["David", "Frank"], ["Emma", "Grace"],
  ["Frank", "Henry"], ["Grace", "Ivy"], ["Henry", "Jack"],
  ["Ivy", "Kate"], ["Jack", "Liam"], ["Kate", "Mia"],
  ["Liam", "Noah"], ["Mia", "Olivia"], ["Noah", "Paul"],
  ["Olivia", "Quinn"], ["Paul", "Ruby"], ["Bob", "Emma"],
  ["Charlie", "Frank"], ["David", "Grace"], ["Emma", "Henry"],
  ["Frank", "Ivy"], ["Grace", "Jack"], ["Henry", "Kate"]
];

const FAKE_SEARCH_PATH = {
  from: "Alice",
  to: "Ruby",
  path: ["Alice", "Bob", "David", "Frank", "Henry", "Jack", "Liam", "Noah", "Paul", "Ruby"],
  explored: ["Alice", "Bob", "Charlie", "David", "Emma", "Frank", "Grace", "Henry", "Ivy", "Jack", "Kate", "Liam", "Mia", "Noah", "Olivia", "Paul", "Quinn", "Ruby"]
};

const FAKE_SEARCH_LOGS = [
  { id: "1", from: "Alice", to: "Ruby", pathLength: 9, exploredCount: 18, timestamp: "10:23:45 AM" },
  { id: "2", from: "Bob", to: "Kate", pathLength: 5, exploredCount: 12, timestamp: "10:18:32 AM" },
  { id: "3", from: "Charlie", to: "Jack", pathLength: 4, exploredCount: 10, timestamp: "10:12:18 AM" },
  { id: "4", from: "Emma", to: "Paul", pathLength: 6, exploredCount: 14, timestamp: "10:05:27 AM" },
  { id: "5", from: "Grace", to: "Noah", pathLength: 3, exploredCount: 8, timestamp: "09:58:13 AM" },
];

const FAKE_NODE_RANKINGS = [
  { name: "Alice", count: 15 },
  { name: "Bob", count: 12 },
  { name: "Emma", count: 10 },
  { name: "Jack", count: 9 },
  { name: "Kate", count: 8 },
  { name: "Charlie", count: 7 },
  { name: "Paul", count: 6 },
  { name: "Ruby", count: 5 },
];

export default function Home() {
  const [fromNode, setFromNode] = useState("");
  const [toNode, setToNode] = useState("");
  const [hoveredNode, setHoveredNode] = useState<string>("");
  const [hasSearched, setHasSearched] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const sigmaRef = useRef<any>(null);
  const graphRef = useRef<Graph>(new Graph());

  useEffect(() => {
    const graph = graphRef.current;

    FAKE_NODES.forEach((node) => {
      if (!graph.hasNode(node)) {
        graph.addNode(node, { label: node, size: 10, color: "#475569" });
      }
    });

    FAKE_EDGES.forEach(([source, target]) => {
      if (!graph.hasEdge(source, target)) {
        graph.addEdge(source, target, { size: 2, color: "#334155" });
      }
    });

    circular.assign(graph);
    initializeSigma();

    return () => {
      if (sigmaRef.current) {
        sigmaRef.current.kill();
      }
    };
  }, []);

  const initializeSigma = async () => {
    if (!containerRef.current) return;

    const Sigma = (await import("sigma")).default;

    if (sigmaRef.current) {
      sigmaRef.current.kill();
    }

    const renderer = new Sigma(graphRef.current, containerRef.current, {
      renderEdgeLabels: false,
      defaultNodeColor: "#475569",
      defaultEdgeColor: "#334155",
      labelColor: { color: "#e2e8f0" },
      labelSize: 12,
    });

    renderer.on("enterNode", ({ node }) => {
      setHoveredNode(node);
    });

    renderer.on("leaveNode", () => {
      setHoveredNode("");
    });

    sigmaRef.current = renderer;
  };

  const handleSearch = () => {
    if (!fromNode.trim() || !toNode.trim()) return;

    setHasSearched(true);
    visualizeFakePath();
  };

  const visualizeFakePath = () => {
    const graph = graphRef.current;
    const { path, explored } = FAKE_SEARCH_PATH;

    graph.forEachNode((node) => {
      let color = "#475569";
      let size = 10;

      if (explored.includes(node)) {
        color = "#64748b";
        size = 12;
      }

      if (path.includes(node)) {
        color = "#06b6d4";
        size = 15;
      }

      if (node === path[0]) {
        color = "#10b981";
        size = 18;
      }

      if (node === path[path.length - 1]) {
        color = "#f59e0b";
        size = 18;
      }

      graph.setNodeAttribute(node, "color", color);
      graph.setNodeAttribute(node, "size", size);
    });

    graph.forEachEdge((edge) => {
      const [source, target] = graph.extremities(edge);
      const sourceIndex = path.indexOf(source);
      const targetIndex = path.indexOf(target);

      if (
        sourceIndex !== -1 &&
        targetIndex !== -1 &&
        Math.abs(sourceIndex - targetIndex) === 1
      ) {
        graph.setEdgeAttribute(edge, "color", "#06b6d4");
        graph.setEdgeAttribute(edge, "size", 4);
      } else {
        graph.setEdgeAttribute(edge, "color", "#334155");
        graph.setEdgeAttribute(edge, "size", 2);
      }
    });

    if (sigmaRef.current) {
      sigmaRef.current.refresh();
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <div className="container mx-auto px-4 py-8">
        <div className="text-center mb-8">
          <h1 className="text-5xl font-bold mb-2 bg-gradient-to-r from-cyan-400 to-blue-500 bg-clip-text text-transparent">
            6MEET
          </h1>
          <p className="text-slate-400 text-sm">Six Degrees of Separation Visualizer</p>
        </div>

        <Card className="bg-slate-900 border-slate-800 mb-8">
          <CardContent className="pt-6">
            <div className="flex flex-col md:flex-row gap-4 items-end">
              <div className="flex-1">
                <label className="block text-sm font-medium mb-2 text-slate-300">From</label>
                <Input
                  value={fromNode}
                  onChange={(e) => setFromNode(e.target.value)}
                  placeholder="e.g., Alice"
                  className="bg-slate-800 border-slate-700 text-slate-100 placeholder:text-slate-500"
                />
              </div>
              <div className="flex-1">
                <label className="block text-sm font-medium mb-2 text-slate-300">To</label>
                <Input
                  value={toNode}
                  onChange={(e) => setToNode(e.target.value)}
                  placeholder="e.g., Ruby"
                  className="bg-slate-800 border-slate-700 text-slate-100 placeholder:text-slate-500"
                />
              </div>
              <Button
                onClick={handleSearch}
                className="bg-cyan-600 hover:bg-cyan-700 text-white px-8"
              >
                <Search className="w-4 h-4 mr-2" />
                Search
              </Button>
            </div>
            {!hasSearched && (
              <p className="text-sm text-slate-500 mt-4 text-center">
                Try searching from "Alice" to "Ruby" to visualize the path
              </p>
            )}
          </CardContent>
        </Card>

        <Card className="bg-slate-900 border-slate-800 mb-8">
          <CardHeader>
            <CardTitle className="text-cyan-400 flex justify-between items-center">
              <span>Graph Visualization</span>
              {hasSearched && (
                <span className="text-sm font-normal text-slate-400">
                  Explored Nodes: {FAKE_SEARCH_PATH.explored.length}
                </span>
              )}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div
              ref={containerRef}
              className="w-full h-[500px] bg-slate-950 rounded-lg relative border border-slate-800"
            >
              {hoveredNode && (
                <div className="absolute top-4 left-4 bg-slate-800 px-4 py-2 rounded-lg border border-cyan-500 shadow-lg z-10">
                  <p className="text-cyan-400 font-medium">{hoveredNode}</p>
                </div>
              )}
            </div>
            <div className="mt-4 flex flex-wrap gap-4 text-sm">
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-green-500"></div>
                <span className="text-slate-300">Start Node</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-amber-500"></div>
                <span className="text-slate-300">End Node</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-cyan-500"></div>
                <span className="text-slate-300">Path</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-slate-500"></div>
                <span className="text-slate-300">Explored</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className="grid md:grid-cols-2 gap-8 mb-8">
          <Card className="bg-slate-900 border-slate-800">
            <CardHeader>
              <CardTitle className="text-cyan-400">Search Log</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow className="border-slate-800">
                      <TableHead className="text-slate-400">From</TableHead>
                      <TableHead className="text-slate-400">To</TableHead>
                      <TableHead className="text-slate-400">Hops</TableHead>
                      <TableHead className="text-slate-400">Explored</TableHead>
                      <TableHead className="text-slate-400">Time</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {FAKE_SEARCH_LOGS.map((log) => (
                      <TableRow key={log.id} className="border-slate-800">
                        <TableCell className="text-slate-300">{log.from}</TableCell>
                        <TableCell className="text-slate-300">{log.to}</TableCell>
                        <TableCell className="text-slate-300">{log.pathLength}</TableCell>
                        <TableCell className="text-slate-300">{log.exploredCount}</TableCell>
                        <TableCell className="text-slate-300">{log.timestamp}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>

          <Card className="bg-slate-900 border-slate-800">
            <CardHeader>
              <CardTitle className="text-cyan-400">Node Rankings</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow className="border-slate-800">
                      <TableHead className="text-slate-400">Rank</TableHead>
                      <TableHead className="text-slate-400">Node</TableHead>
                      <TableHead className="text-slate-400">Searches</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {FAKE_NODE_RANKINGS.map((rank, index) => (
                      <TableRow key={rank.name} className="border-slate-800">
                        <TableCell className="text-slate-300">{index + 1}</TableCell>
                        <TableCell className="text-slate-300 font-medium">{rank.name}</TableCell>
                        <TableCell className="text-slate-300">{rank.count}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>

        <Card className="bg-slate-900 border-slate-800 mb-8">
          <CardHeader>
            <CardTitle className="text-cyan-400">Six Degrees of Separation</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-slate-300 leading-relaxed">
              The theory of <strong>Six Degrees of Separation</strong> suggests that any two people on Earth are connected
              through a chain of six or fewer social connections. This concept demonstrates the "small world" phenomenon,
              where despite the vast number of people in the world, we are all surprisingly interconnected through our
              social networks. Each person is connected to every other person by an average of six steps or fewer, forming
              a web of relationships that spans the globe. This visualization demonstrates how these connection chains
              work, showing the path between individuals and all the nodes explored in finding that connection.
            </p>
          </CardContent>
        </Card>

        <div className="text-center pb-8">
          <a
            href="https://github.com"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 text-slate-400 hover:text-cyan-400 transition-colors"
          >
            <Github className="w-5 h-5" />
            <span>View on GitHub</span>
          </a>
        </div>
      </div>
    </div>
  );
}
