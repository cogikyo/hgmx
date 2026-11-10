import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const rows = el.querySelectorAll('path[d$="h18"]');
  animate(rows, { scaleX: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
  const cut = el.querySelector('path[d="M12 15V9"]');
  animate(cut, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
