import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = [el.querySelector("rect"), el.querySelector('path[d="M9 3v18"]')];
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const chevron = el.querySelector('path[d^="m16"]');
  animate(chevron, { translateX: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
