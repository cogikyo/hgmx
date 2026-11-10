import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const chevron = el.querySelector('path[d^="m3 8"]');
  animate(chevron, { translateX: [0, 2, -0.5, 0] }, { duration: τmid });
  const lines = el.querySelectorAll('path[d$="H11"]');
  animate(lines, { translateX: [0, 2, -0.5, 0] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
