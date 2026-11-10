import { animate, τmax } from "../motion.js";

export default function (el) {
  const sides = [
    ["M3 6v12", { translateX: [0, -2, 0.5, 0] }],
    ["M21 6v12", { translateX: [0, 2, -0.5, 0] }],
    ["M6 3h12", { translateY: [0, -2, 0.5, 0] }],
    ["M6 21h12", { translateY: [0, 2, -0.5, 0] }],
  ];
  for (const [d, nudge] of sides) animate(el.querySelector(`path[d="${d}"]`), nudge, { duration: τmax });
}
