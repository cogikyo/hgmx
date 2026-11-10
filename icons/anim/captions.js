import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lines = el.querySelector("path");
  animate(lines, { opacity: [0, 1], x: [-1.5, 0] }, { duration: τmid, delay: τmin });
}
