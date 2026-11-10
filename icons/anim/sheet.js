import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const rows = el.querySelectorAll('line[x1="3"]');
  animate(rows, { scaleX: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
  const columns = el.querySelectorAll('line[y1="9"][y2="21"]');
  animate(columns, { scaleY: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
