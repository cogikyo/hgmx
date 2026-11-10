import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const pad = el.querySelector("rect");
  animate(pad, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
  const lines = ["M8 10h6", "M8 14h8", "M8 18h5"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { opacity: [0, 1], translateX: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
