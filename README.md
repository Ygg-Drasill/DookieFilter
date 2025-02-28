# DookieFilter

```mermaid
graph TD;
    DataStream((("Data")));
    PB[(PringleBuffer)];
    MinIO[(MinIO)];
    Visunator[[Visunator]]

    DataStream ==>|Raw Data|Collector;
    PB -->MinIO & Visunator;

    subgraph VM
    Collector(["Collector"]);
    Detector["Detector"];
    FixSwap["Fix Swap"];
    FixJump["Fix Jump"];
    Filter(("Filter"));

    Collector -->Detector;
    Detector -.->|Swap| FixSwap;
    Detector -.->|Jump| FixJump;
    FixSwap & FixJump -->Filter;
    end
    Collector & Filter -->PB;

    Detector -.->|Missing|AIInterference;
    AIInterference -.-|Read|PB;

    subgraph GPU
    AIInterference["Inference"];
    end

    AIInterference ==o|Predicted Data|Collector;
```
