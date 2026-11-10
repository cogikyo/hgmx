import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const trace = el.querySelector("path");
  animate(trace, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
  animate(trace, { scaleY: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
