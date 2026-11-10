import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const trace = el.querySelector("path");
  animate(trace, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(trace, { scaleY: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
