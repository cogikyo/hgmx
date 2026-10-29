import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const header = el.querySelector('rect[width="18"]');
  animate(header, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const blocks = el.querySelectorAll('rect[y="14"]');
  animate(blocks, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
