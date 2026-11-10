import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const rows = el.querySelectorAll('path[d$="H3"]');
  animate(rows, { scaleX: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
  const column = el.querySelector('path[d="M15 3v18"]');
  animate(column, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
