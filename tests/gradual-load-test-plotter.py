import json
import sys

import matplotlib.pyplot as plt
import pandas as pd


def plot(file_name):
    points = []

    print(f"Loading data from {file_name}...")
    with open(file_name, "r") as f:
        for line in f:
            entry = json.loads(line)
            if entry["type"] == "Point":
                data = entry["data"]
                tags = data.get("tags", {})

                points.append(
                    {
                        "time": pd.to_datetime(data["time"]),
                        "metric": entry["metric"],
                        "value": data["value"],
                        "expected_response": tags.get("expected_response") == "true",
                    }
                )

    df = pd.DataFrame(points)

    # Determine seconds since the start of the test
    start_time = df["time"].min()
    df["seconds"] = (df["time"] - start_time).dt.total_seconds()

    fig, ax_latency = plt.subplots(figsize=(10, 5))

    # Plot successful requests only
    success_latency = df[
        (df["metric"] == "http_req_duration") & df["expected_response"]
    ]

    ax_latency.scatter(
        success_latency["seconds"],
        success_latency["value"],
        alpha=0.4,
        s=8,
        color="#3498db",
        label="Successful Request Latency (ms)",
    )

    ax_latency.set_ylabel("Latency (ms)", fontweight="bold")
    ax_latency.set_xlabel("Time into test (s)", fontweight="bold")

    ax_latency.set_ylim(bottom=1, top=500)
    ax_latency.set_xlim(0, 60)

    ax_latency.grid(True, which="both", ls="-", alpha=0.2)

    ax_rps = ax_latency.twinx()

    # Plot requests per second (RPS)
    # I do this by resampling the data to 1 second intervals and counting the number of iterations
    # for each interval
    rps_df = df[df["metric"] == "iterations"].set_index("time")
    rps_series = rps_df["value"].resample("1s").count()

    rps_x = range(len(rps_series))

    ax_rps.plot(
        rps_x,
        rps_series,
        color="#e67e22",
        linewidth=2,
    )

    ax_rps.set_ylabel("Requests per second (RPS)", fontweight="bold", color="#e67e22")
    ax_rps.tick_params(axis="y", labelcolor="#e67e22")

    ax_rps.set_ylim(bottom=0, top=100)

    ax_latency.legend(
        loc="upper left",
        frameon=True,
        shadow=True,
    )

    plt.title("Latency vs. RPS in a Gradual Load Test", fontsize=14)
    plt.tight_layout()
    plt.show()


if __name__ == "__main__":
    if len(sys.argv) > 1:
        plot(sys.argv[1])
    else:
        print("Usage: python script.py results.json")
