import React, { useMemo, useState, useRef } from "react";

import { motion, AnimatePresence } from "framer-motion";
import styles from "./piecharts.module.scss";

// Default export: PieChart component
// Usage example at the bottom of the file.

// Types
export type PieDatum = {
  id: string | number;
  label: string;
  value: number;
  color?: string; // optional custom color
};

export type LegendPosition = "right" | "bottom" | "left" | "top";

export type PieChartProps = {
  data: PieDatum[];
  size?: number; // width & height of the SVG in px (square)
  innerRadius?: number; // 0 = pie, >0 = donut
  padAngleDeg?: number; // degrees gap between slices
  colors?: string[]; // fallback palette
  showLegend?: boolean;
  legendPosition?: LegendPosition;
  animate?: boolean;
  className?: string;
  onSliceClick?: (d: PieDatum) => void;
  tooltipFormatter?: (d: PieDatum) => React.ReactNode;
};

// Helpers
const TAU = Math.PI * 2;

function polarToCartesian(cx: number, cy: number, r: number, angleDeg: number) {
  const angleRad = ((angleDeg - 90) * Math.PI) / 180.0; // start at 12 o'clock
  return {
    x: cx + r * Math.cos(angleRad),
    y: cy + r * Math.sin(angleRad),
  };
}

function describeArc(
  cx: number,
  cy: number,
  rOuter: number,
  rInner: number,
  startAngle: number,
  endAngle: number
) {
  // Guard for full circle
  const full = Math.abs(endAngle - startAngle) >= 360 - 1e-6;
  if (full) {
    // Draw full circle as two semicircles (SVG arc can't draw full circle in one arc)
    const top = `M ${cx} ${cy - rOuter} A ${rOuter} ${rOuter} 0 1 1 ${cx - 0.0001} ${cy - rOuter}`;
    if (rInner <= 0) return `${top} Z`;
    // donut full circle
    const inner = `M ${cx} ${cy - rInner} A ${rInner} ${rInner} 0 1 0 ${cx - 0.0001} ${cy - rInner}`;
    return `${top} ${inner} Z`;
  }

  const startOuter = polarToCartesian(cx, cy, rOuter, endAngle);
  const endOuter = polarToCartesian(cx, cy, rOuter, startAngle);
  const startInner = polarToCartesian(cx, cy, rInner, startAngle);
  const endInner = polarToCartesian(cx, cy, rInner, endAngle);

  const largeArcFlag = endAngle - startAngle <= 180 ? "0" : "1";

  const d = [
    `M ${startOuter.x} ${startOuter.y}`,
    `A ${rOuter} ${rOuter} 0 ${largeArcFlag} 0 ${endOuter.x} ${endOuter.y}`,
    `L ${startInner.x} ${startInner.y}`,
    rInner > 0
      ? `A ${rInner} ${rInner} 0 ${largeArcFlag} 1 ${endInner.x} ${endInner.y}`
      : `L ${cx} ${cy}`,
    "Z",
  ].join(" ");

  return d;
}

const DEFAULT_PALETTE = [
  "#F59E0B",
  "#06B6D4",
  "#10B981",
  "#EF4444",
  "#4F46E5",
  "#8B5CF6",
  "#F97316",
  "#84CC16",
];

// Main component
export default function PieChart({
  data,
  size = 340,
  innerRadius = 0,
  padAngleDeg = 0,
  colors = DEFAULT_PALETTE,
  showLegend = true,
  animate = true,
  className = "",
  onSliceClick,
  tooltipFormatter,
}: PieChartProps) {
  const total = useMemo(() => data.reduce((s, d) => s + Math.max(0, d.value), 0), [data]);
  const sorted = useMemo(() => [...data].sort((a, b) => b.value - a.value), [data]);
  const palette = colors;

  const radius = size / 2;
  const outerR = radius - 6; // padding
  const innerR = Math.max(0, innerRadius);

  // compute angles
  const slices = useMemo(() => {
    let start = 0;
    const gap = padAngleDeg; // degrees
    const res = sorted.map((d, i) => {
      const frac = total === 0 ? 0 : Math.max(0, d.value) / total;
      const angle = frac * 360;
      const startAngle = start;
      const endAngle = start + angle;
      // apply small padding by reducing endAngle slightly and advancing start a bit
      const pad = gap * (angle / 360); // preserve relative padding
      const adjustedStart = startAngle + pad / 2;
      const adjustedEnd = endAngle - pad / 2;
      start = endAngle;
      return {
        datum: d,
        startAngle: adjustedStart,
        endAngle: adjustedEnd,
        color: d.color || palette[i % palette.length],
        percent: frac * 100,
      };
    });
    return res;
  }, [sorted, total, padAngleDeg, palette]);

  const [hoverId, setHoverId] = useState<string | number | null>(null);
  const [activeId, setActiveId] = useState<string | number | null>(null);

  const tooltipRef = useRef<HTMLDivElement | null>(null);

  // Accessible keyboard handler
  function handleKeyPress(e: React.KeyboardEvent, d: PieDatum) {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      onSliceClick?.(d);
      setActiveId(d.id);
    }
  }

  // Layout helpers for legend


  // Tooltip content
  const renderTooltip = (d: PieDatum | null) => {
    if (!d) return null;
    if (tooltipFormatter) return tooltipFormatter(d);
    return (
      <div className="text-xs">
        <div className="font-medium">{d.label}</div>
        <div className="text-neutral-500">{d.value}</div>
      </div>
    );
  };

  return (
    <div className={styles.containerClass}>
      <div className={styles.chartContainer}>
        <svg
          width={size}
          height={size}
          viewBox={`0 0 ${size} ${size}`}
          role="img"
          aria-label="Pie chart"
          className="block">
          <title>Pie Chart</title>
          <g transform={`translate(${radius}, ${radius})`}>
            {/* slices */}
            {slices.map((s, i) => {
              const d = describeArc(0, 0, outerR, innerR, s.startAngle, s.endAngle);
              const hover = hoverId === s.datum.id;
              const active = activeId === s.datum.id;
              const offset = hover || active ? 8 : 0;

              // compute centroid for label / a small translate for pop-out effect
              const mid = (s.startAngle + s.endAngle) / 2;
              const tx = (offset * Math.cos(((mid - 90) * Math.PI) / 180)).toFixed(3);
              const ty = (offset * Math.sin(((mid - 90) * Math.PI) / 180)).toFixed(3);

              return (
                <g
                  key={s.datum.id}
                  transform={`translate(${tx}, ${ty})`}
                  /*onMouseEnter={() => setHoverId(s.datum.id)}
                  onMouseLeave={() => setHoverId(null)}*/
                >
                  <AnimatePresence>
                    <motion.path
                      d={d}
                      fill={s.color}
                      stroke="#ffffff"
                      strokeWidth={1}
                      initial={{ opacity: 0, pathLength: 0 }}
                      animate={{ opacity: 1, pathLength: 1 }}
                      exit={{ opacity: 0, pathLength: 0 }}
                      transition={{ duration: animate ? 0.6 : 0 }}
                      //style={{ cursor: onSliceClick ? "pointer" : "default" }}
                     // onClick={() => onSliceClick?.(s.datum)}
                      role="button"
                      tabIndex={0}
                      onKeyDown={(e) => handleKeyPress(e, s.datum)}
                      aria-label={`${s.datum.label}: ${s.datum.value} (${s.percent.toFixed(1)}%)`}
                    />
                  </AnimatePresence>

                  {/* optional small label near outer arc when hovered */}
                  {/*hover && (
                    <text
                      x={
                        (outerR + 12) * Math.cos(((mid - 90) * Math.PI) / 180)
                      }
                      y={
                        (outerR + 12) * Math.sin(((mid - 90) * Math.PI) / 180)
                      }
                      textAnchor={mid > 90 && mid < 270 ? "end" : "start"}
                      alignmentBaseline="middle"
                      fontSize={12}
                      className="font-medium"
                    >
                      {s.datum.label} — {s.percent.toFixed(1)}%
                    </text>
                  )*/}
                </g>
              );
            })}

            {/* center label: total */}
            <g>
              <text
                x={0}
                y={0}
                textAnchor="middle"
                alignmentBaseline="middle"
                className="text-sm font-semibold"
                style={{ fontSize: 14 }}
              >
                
              </text>
            </g>
          </g>
        </svg>

        {/* simple tooltip */}
        <div
          ref={tooltipRef}
          role="status"
          aria-live="polite"
          className="pointer-events-none absolute left-1/2 top-0 -translate-x-1/2 mt-1"
        >
          {hoverId && (
            <div className="bg-white border rounded shadow px-2 py-1 text-xs">{
              renderTooltip(data.find((x) => x.id === hoverId) ?? null)
            }</div>
          )}
        </div>
      </div>

    
{showLegend && (
  <ul className={styles.legendList}>
    {slices.map((s) => (
      <li key={s.datum.id} className={styles.legendItem}>
        <span
          className={styles.legendColor}
          style={{ background: s.color, display: "inline-block" }}
        />
        <span className={styles.legendLabel}>{s.datum.label}</span>
      </li>
    ))}
  </ul>
)}
    </div>
  );
}


/* ------------------ Example usage ------------------

import React from 'react';
import PieChart, { PieDatum } from './react-piechart-library';

const data: PieDatum[] = [
  { id: 'a', label: 'Apples', value: 40 },
  { id: 'b', label: 'Bananas', value: 25 },
  { id: 'c', label: 'Cherries', value: 20 },
  { id: 'd', label: 'Dates', value: 15 },
];

export default function Demo() {
  return (
    <div className="p-6">
      <h2 className="text-xl mb-4">Sales by fruit</h2>
      <PieChart
        data={data}
        size={420}
        innerRadius={80}
        padAngleDeg={1}
        showLegend={true}
        legendPosition="right"
        onSliceClick={(d) => alert(`clicked ${d.label}`)}
      />
    </div>
  );
}

----------------------------------------------------

Notes / next steps ideas you may request:
- Add label placement with collision detection
- Export as SVG / PNG (render offscreen and toBlob)
- Add accessibility improvements (aria-describedby for tooltip, focus ring styling)
- Add unit tests and storybook stories
- Add bundling: rollup/tsup config and published npm package with types
- Add small hook usePieLayout for SSR-friendly calculations
*/
