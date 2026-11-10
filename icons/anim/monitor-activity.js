import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector("rect");
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const trace = el.querySelector('path[d^="M18 10"]');
  animate(trace, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(trace, { scaleY: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
