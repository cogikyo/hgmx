import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector("rect");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const trace = el.querySelector('path[d^="M6 12"]');
  animate(trace, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(trace, { scaleY: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
