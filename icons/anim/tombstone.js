import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const stone = el.querySelector('path[d^="M19 21V10"]');
  animate(stone, { scaleY: [0, 1.1, 1] }, { duration: τmid });
  const lines = el.querySelectorAll('path[d="M11 10h2"], path[d="M9 14h6"]');
  animate(lines, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
