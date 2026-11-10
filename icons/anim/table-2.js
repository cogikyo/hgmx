import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const row = el.querySelector('path[d="M3 9h18"]');
  animate(row, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const column = el.querySelector('path[d="M9 3v18"]');
  animate(column, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
