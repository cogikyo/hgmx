import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const frame = ["M4 8V6", "M16 4h2", "M20 12v2", "M20 18v2", "M13 22h-2", "M7 22H6", "M4 14v-2"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });
  const lines = ["M8 10h6", "M8 14h8", "M8 18h5"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { opacity: [0, 1], translateX: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
